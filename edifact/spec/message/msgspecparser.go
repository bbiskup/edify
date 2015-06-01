package message

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const pathSep = string(os.PathSeparator)

var (
	// e.g.
	// "SOURCE: TBG3 Transport"
	sourceRE = regexp.MustCompile(`^SOURCE: (.*) *$`)

	// e.g.
	// "00210       ---- Segment group 6  ------------------ C   99-------------+||"
	segmentGroupStartRE = regexp.MustCompile(
		`^(\d{5})[ *]+-{4} Segment group (\d+)[ ]*[-]+ ([MC])[ ]*(\d+)[ ]*[-]+[+]+([|]*)$`)

	// e.g. (top-level; not in group)
	//"00010   UNH Message header                           M   1     "
	//
	// in group
	// "00060   TDT Transport information                    M   1                |"
	//
	// in group; deeper nesting
	// "00140   DTM Date/time/period                         M   1               ||"
	//
	// group end (2 nesting levels at once)
	// "00160   QTY Quantity                                 C   99--------------++"
	segmentRE = regexp.MustCompile(`^(\d{5})[ ]{3}([A-Z]{3}) (.{20,}) ([MC])[ ]*(\d+)[-+ ]*([|]*)[ ]*$`)

	// A segment spec that spans multiple lines because
	// of a long name, e.g. (QALITY_D.14B)
	// "00250   SPS Sampling parameters for summary                               |"
	// "               statistics                            C   1                |"
)

type SegmentGroupStart struct {
	RecordNum    int
	GroupNum     int
	IsMandatory  bool
	MaxCount     int
	NestingLevel int
}

type SegmentEntry struct {
	RecordNum   int
	SegmentId   string
	SegmentName string
	IsMandatory bool
	MaxCount    int

	// Nesting level _after_ this segment entry
	// A segment entry might close multiple groups simultaneously.
	NestingLevel int
}

// Used for parallel parsing of segment specs
type FileSpec struct {
	fileName string
	contents string
}

// Parser for message specifications
// e.g. d14b/edmd/AUTHOR_D.14B
type MsgSpecParser struct {
	segmentSpecs segment.SegmentSpecProvider
}

func (p *MsgSpecParser) parseDate(dateStr string) (date time.Time, err error) {
	date, err = time.Parse("2006-01-02", dateStr)
	return
}

func (p *MsgSpecParser) parseSource(sourceStr string) (source string, err error) {
	match := sourceRE.FindStringSubmatch(sourceStr)
	if match == nil {
		return "", errors.New(fmt.Sprintf("Could not get source from '%s'",
			sourceStr))
	}

	if len(match) != 2 {
		panic("Internal error: incorrect regular expression")
	}

	return match[1], nil
}

// lines: lines of message spec file (without header)
func (p *MsgSpecParser) getSegmentTableLines(lines []string) (segmentTable []string, err error) {
	started := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Pos     Tag Name") {
			started = true
			continue
		}
		if started {
			segmentTable = append(segmentTable, line)
		}
	}
	if !started {
		err = errors.New(fmt.Sprintf("Segment table not found"))
	}
	return
}

func (p *MsgSpecParser) logNestingLevelChange(currentNestingLevel int, newNestingLevel int) {
	if currentNestingLevel != newNestingLevel {
		// log.Printf("  Nesting level %d ---> %d",
		//	currentNestingLevel, newNestingLevel)
	} else {
		// log.Printf("  Nesting level remains at %d", currentNestingLevel)
	}
}

func (p *MsgSpecParser) shouldSkipSegTableLine(line string) bool {
	if len(strings.TrimSpace(line)) == 0 {
		return true
	}

	if p.matchHeaderOrEmptyInGroupSection(line) {
		return true
	}
	return false
}

// Join multi-line segment definition
func (p *MsgSpecParser) joinMultiLineSegmentDef(
	line string, index int,
	numLines int, segmentTableLines []string) (joinedLine string, newIndex int) {

	if index < numLines-1 {
		nextLine := segmentTableLines[index+1]
		if strings.HasPrefix(nextLine, "               ") && !strings.HasPrefix(nextLine, "                      ") {
			// log.Printf("Joining multi-line")
			firstLine := strings.TrimRight(line, "+| ")
			line = firstLine + " " + strings.TrimSpace(nextLine)
			index++
		}
	}
	return line, index
}

/*
 Get sequence of segments/segment groups
 e.g.

00010   UNH Message header                           M   1
00020   BGM Beginning of message                     M   1
00030   DTM Date/time/period                         C   1
00040   BUS Business function                        C   1

00050       ---- Segment group 1  ------------------ C   2----------------+
00060   RFF Reference                                M   1                |
00070   DTM Date/time/period                         C   1----------------+

00080       ---- Segment group 2  ------------------ C   5----------------+
00090   FII Financial institution information        M   1                |
00100   CTA Contact information                      C   1                |
00110   COM Communication contact                    C   5----------------+
...
*/
func (p *MsgSpecParser) parseMsgSpecParts(fileName string, lines []string) (msgSpecParts []MsgSpecPart, err error) {
	segmentTableLines, err := p.getSegmentTableLines(lines)
	currentNestingLevel := 0
	var currentMsgSpecPart MsgSpecPart = nil
	numLines := len(segmentTableLines)

	for index, line := range segmentTableLines {
		line = strings.TrimRight(line, " \r\n")
		if p.shouldSkipSegTableLine(line) {
			continue
		}

		line, index := p.joinMultiLineSegmentDef(line, index, numLines, segmentTableLines)

		// Each line must either be a segment entry or segment group start
		segmentEntry, err := p.parseSegmentEntry(line)
		if err != nil {
			return nil, err
		}

		// Are we dealing with a segment entry?
		if segmentEntry != nil {
			p.logNestingLevelChange(currentNestingLevel, segmentEntry.NestingLevel)
			nestingDelta := currentNestingLevel - segmentEntry.NestingLevel

			segmentSpec := p.segmentSpecs.Get(segmentEntry.SegmentId)
			if segmentSpec == nil {
				return nil, errors.New(fmt.Sprintf("No segment spec for ID '%s'",
					segmentEntry.SegmentId))
			}
			part := NewMsgSpecSegmentPart(
				segmentSpec, segmentEntry.MaxCount, segmentEntry.IsMandatory, currentMsgSpecPart)

			if currentNestingLevel == 0 {
				msgSpecParts = append(msgSpecParts, part)
			} else {
				group, ok := currentMsgSpecPart.(*MsgSpecSegmentGroupPart)
				if !ok {
					return nil, errors.New(fmt.Sprintf(
						"Internal error: nesting incorrect; got: %#v",
						currentMsgSpecPart))
				}
				group.Append(part)
			}

			if nestingDelta > 0 {
				// Navigate up in message spec part hierarchy
				for level := 0; level < nestingDelta; level++ {
					currentMsgSpecPart = currentMsgSpecPart.Parent()
				}
			}

			currentNestingLevel = segmentEntry.NestingLevel
			continue
		}

		// Next alternative: segment group start
		segmentGroupStartSpec, err := p.parseSegmentGroupStart(line)
		if err != nil {
			return nil, err
		}

		sg := segmentGroupStartSpec
		if sg != nil {
			p.logNestingLevelChange(currentNestingLevel, sg.NestingLevel)

			group := NewMsgSpecSegmentGroupPart(
				fmt.Sprintf("Group_%d", sg.GroupNum),
				[]MsgSpecPart{}, sg.MaxCount,
				sg.IsMandatory, currentMsgSpecPart)

			if currentMsgSpecPart == nil {
				msgSpecParts = append(msgSpecParts, group)
			} else {
				parentGroup, ok := currentMsgSpecPart.(*MsgSpecSegmentGroupPart)
				if !ok {
					return nil, errors.New(fmt.Sprintf(
						"Internal error: nesting incorrect; got: %#v",
						currentMsgSpecPart))
				}
				parentGroup.Append(group)
			}
			currentMsgSpecPart = group
			currentNestingLevel = sg.NestingLevel
		} else {
			return nil, errors.New(
				fmt.Sprintf("Parse error in file %s at index %d ('%s')",
					fileName, index, line))
		}
	}
	return
}

func (p *MsgSpecParser) parseSegmentGroupStart(line string) (segmentGroupStart *SegmentGroupStart, err error) {
	match := segmentGroupStartRE.FindStringSubmatch(line)
	if match == nil {
		// not an error; other pattern might still match
		return
	}

	if len(match) != 6 {
		panic("Internal error: incorrect regular expression")
	}

	recordNum, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	groupNum, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	isMandatoryStr := match[3]
	var isMandatory bool
	switch isMandatoryStr {
	case "C":
		isMandatory = false
	case "M":
		isMandatory = true
	}

	maxCount, err := strconv.Atoi(match[4])
	if err != nil {
		return
	}

	bars := match[5]

	return &SegmentGroupStart{
		RecordNum:   recordNum,
		GroupNum:    groupNum,
		IsMandatory: isMandatory,
		MaxCount:    maxCount,

		// Conceptually, nesting level 0 is outside any group
		NestingLevel: len(bars) + 1,
	}, nil
}

func (p *MsgSpecParser) parseSegmentEntry(line string) (segmentEntry *SegmentEntry, err error) {
	match := segmentRE.FindStringSubmatch(line)
	if match == nil {
		return
	}

	if len(match) != 7 {
		panic("Internal error: incorrect regular expression")
	}

	recordNum, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	isMandatoryStr := match[4]
	var isMandatory bool
	switch isMandatoryStr {
	case "C":
		isMandatory = false
	case "M":
		isMandatory = true
	}

	maxCount, err := strconv.Atoi(match[5])
	if err != nil {
		return
	}

	bars := match[6]

	return &SegmentEntry{
		RecordNum:    recordNum,
		SegmentId:    match[2],
		SegmentName:  strings.TrimSpace(match[3]),
		IsMandatory:  isMandatory,
		MaxCount:     maxCount,
		NestingLevel: len(bars),
	}, nil
}

// e.g.
// "            HEADER SECTION"
// or
// "
func (p *MsgSpecParser) matchHeaderOrEmptyInGroupSection(line string) (matches bool) {
	return strings.HasPrefix(line, "            ")
}

func (p *MsgSpecParser) getFileContents(fileName string) (contents []byte, err error) {
	return ioutil.ReadFile(fileName)
}

func (p *MsgSpecParser) ParseSpecFile(fileName string) (spec *MsgSpec, err error) {
	contents, err := p.getFileContents(fileName)
	if err != nil {
		return nil, err
	}
	return p.ParseSpecFileContents(fileName, string(contents))
}

// One spec file contains the spec for a single message type
/*
                                UN/EDIFACT

                  UNITED NATIONS STANDARD MESSAGE (UNSM)

                              Invoice message
...
                                           Message Type : INVOIC
                                           Version      : D
                                           Release      : 14B
                                           Contr. Agency: UN

                                           Revision     : 16
                                           Date         : 2014-11-17
...
SOURCE: TBG1 Supply Chain
*/
func (p *MsgSpecParser) ParseSpecFileContents(fileName string, contents string) (spec *MsgSpec, err error) {
	// The largest standard message file has 321k (about 6800 lines), so
	// we can read it at once

	// log.Printf("Parsing message spec file contents'%s'", fileName)

	lines := strings.Split(contents, "\n")
	minLines := 48
	numLines := len(lines)
	if numLines < minLines {
		return nil, errors.New(
			fmt.Sprintf("Spec file does not contain enough lines (got %d, expected %d)",
				numLines, minLines))
	}

	name := strings.TrimSpace(lines[4])
	fmt.Printf("id-line: '%s'", lines[33])

	detailCol := 58

	id := strings.TrimSpace(lines[33][detailCol:])
	version := strings.TrimSpace(lines[34][detailCol:])
	release := strings.TrimSpace(lines[35][detailCol:])
	contrAgency := strings.TrimSpace(lines[36][detailCol:])
	revision := strings.TrimSpace(lines[38][detailCol:])
	date, err := p.parseDate(strings.TrimSpace(lines[39][detailCol:]))
	if err != nil {
		return
	}

	source, err := p.parseSource(lines[46])
	if err != nil {
		return
	}

	specParts, err := p.parseMsgSpecParts(fileName, lines[47:])
	if err != nil {
		return
	}

	return NewMsgSpec(id, name, version, release, contrAgency, revision, date, source, specParts), nil
}

func (p *MsgSpecParser) ParseSpecDir(dirName string, suffix string) (specs []*MsgSpec, err error) {
	return p.parseSpecDir_sequential(dirName, suffix)
}

// Parse segment spec directory sequentially
func (p *MsgSpecParser) parseSpecDir_sequential(dirName string, suffix string) (specs []*MsgSpec, err error) {
	entries, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	specs = []*MsgSpec{}
	for _, entry := range entries {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, "."+suffix) {
			continue
		}

		if strings.HasPrefix(fileName, "EDMDI") {
			// message index file
			continue
		}

		spec, err := p.ParseSpecFile(dirName + pathSep + fileName)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return
}

func (p *MsgSpecParser) parseSpecDir_parallel(
	dirName string, suffix string) (specs []*MsgSpec, err error) {
	fmt.Printf("NumThreads: %d; num go routines %d\n",
		edifact.NumThreads, runtime.NumGoroutine())

	var wg sync.WaitGroup

	entries, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}

	for _, entry := range entries {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, "."+suffix) {
			continue
		}

		if strings.HasPrefix(fileName, "EDMDI") {
			// message index file
			continue
		}
		fileNames = append(fileNames, fileName)
	}

	numFiles := len(fileNames)

	fileSpecCh := make(chan FileSpec, 0)
	resultCh := make(chan *MsgSpec, 0)
	resultsCh := make(chan []*MsgSpec, 1)

	// Collect results
	go func() {
		resultSpecs := []*MsgSpec{}
		for i := 0; i < numFiles; i++ {
			newSpec := <-resultCh
			resultSpecs = append(resultSpecs, newSpec)
		}
		resultsCh <- resultSpecs
		close(resultCh)
	}()

	// Parse in parallel
	for i := 0; i < edifact.NumThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileSpec := range fileSpecCh {
				msgSpec, err := p.ParseSpecFileContents(
					fileSpec.fileName, fileSpec.contents)
				if err != nil {
					panic(fmt.Sprintf("TODO: handle err %s", err))
				}
				resultCh <- msgSpec
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, fileName := range fileNames {
			fullPath := path.Clean(dirName + pathSep + fileName)
			contents, err := p.getFileContents(fullPath)
			if err != nil {
				panic(fmt.Sprintf("TODO: handle err %s", err))
			}
			fileSpec := FileSpec{fullPath, string(contents)}
			fileSpecCh <- fileSpec
		}
		close(fileSpecCh)
	}()

	specs = <-resultsCh
	close(resultsCh)
	wg.Wait()
	return
}

func NewMsgSpecParser(segmentSpecs segment.SegmentSpecProvider) *MsgSpecParser {
	return &MsgSpecParser{
		segmentSpecs: segmentSpecs,
	}
}

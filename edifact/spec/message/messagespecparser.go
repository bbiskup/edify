package message

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Parser for message specifications
// e.g. d14b/edmd/AUTHOR_D.14B
type MessageSpecParser struct {
	segmentSpecs segment.SegmentSpecProvider
}

func (p *MessageSpecParser) parseDate(dateStr string) (date time.Time, err error) {
	date, err = time.Parse("2006-01-02", dateStr)
	return
}

func (p *MessageSpecParser) parseSource(sourceStr string) (source string, err error) {
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
func (p *MessageSpecParser) getSegmentTableLines(lines []string) (segmentTable []string, err error) {
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

func (p *MessageSpecParser) logNestingLevelChange(currentNestingLevel int, newNestingLevel int) {
	if currentNestingLevel != newNestingLevel {
		log.Printf("  Nesting level %d ---> %d",
			currentNestingLevel, newNestingLevel)
	} else {
		// log.Printf("  Nesting level remains at %d", currentNestingLevel)
	}
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
func (p *MessageSpecParser) parseMessageSpecParts(lines []string) (messageSpecParts []MessageSpecPart, err error) {
	segmentTableLines, err := p.getSegmentTableLines(lines)
	log.Printf("######################### -> segmentTableLines")
	log.Printf("%s", strings.Join(segmentTableLines, "\n"))
	log.Printf("######################### <- segmentTableLines")
	currentNestingLevel := 0
	var currentMessageSpecPart MessageSpecPart = nil
	numLines := len(segmentTableLines)
	log.Printf("Processing %d segment table lines", numLines)
	for index, line := range segmentTableLines {
		line = strings.TrimRight(line, " \r\n")
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if p.matchHeaderOrEmptyInGroupSection(line) {
			continue
		}

		// Join multi-line segment definition
		// multiLineStartMatch := segmentMultiLineRE.FindStringSubmatch(line)
		if index < numLines-1 {
			nextLine := segmentTableLines[index+1]
			log.Printf("nextLine '%s'", nextLine)
			if strings.HasPrefix(nextLine, "               ") && !strings.HasPrefix(nextLine, "                      ") {
				log.Printf("Joining multi-line")
				firstLine := strings.TrimRight(line, "+| ")
				line = firstLine + " " + strings.TrimSpace(nextLine)
				index++
			}
		}

		// log.Printf("line %d: '%s'\n", index, line)

		// Each line must either be a segment entry or segment group start

		segmentEntry, err := p.parseSegmentEntry(line)
		if err != nil {
			return nil, err
		}

		// Are we dealing with a segment entry?
		if segmentEntry != nil {
			log.Printf("Found %#v", segmentEntry)
			p.logNestingLevelChange(currentNestingLevel, segmentEntry.NestingLevel)
			nestingDelta := currentNestingLevel - segmentEntry.NestingLevel

			segmentSpec := p.segmentSpecs.Get(segmentEntry.SegmentId)
			if segmentSpec == nil {
				return nil, errors.New(fmt.Sprintf("No segment spec for ID '%s'",
					segmentEntry.SegmentId))
			}
			part := NewMessageSpecSegmentPart(
				segmentSpec, segmentEntry.MaxCount, segmentEntry.IsMandatory, currentMessageSpecPart)

			if currentNestingLevel == 0 {
				messageSpecParts = append(messageSpecParts, part)
			} else {
				group, ok := currentMessageSpecPart.(*MessageSpecSegmentGroupPart)
				if !ok {
					return nil, errors.New(fmt.Sprintf(
						"Internal error: nesting incorrect; got: %#v",
						currentMessageSpecPart))
				}
				group.Append(part)
			}

			if nestingDelta > 0 {
				// Navigate up in message spec part hierarchy
				for level := 0; level < nestingDelta; level++ {
					currentMessageSpecPart = currentMessageSpecPart.Parent()
				}
				log.Printf("Parent now %#v after navigating up %d levels",
					currentMessageSpecPart, nestingDelta)
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
			log.Printf("Found %#v", sg)
			p.logNestingLevelChange(currentNestingLevel, sg.NestingLevel)

			group := NewMessageSpecSegmentGroupPart(
				fmt.Sprintf("Group_%d", sg.GroupNum),
				[]MessageSpecPart{}, sg.MaxCount,
				sg.IsMandatory, currentMessageSpecPart)

			if currentMessageSpecPart == nil {
				messageSpecParts = append(messageSpecParts, group)
			} else {
				parentGroup, ok := currentMessageSpecPart.(*MessageSpecSegmentGroupPart)
				if !ok {
					return nil, errors.New(fmt.Sprintf(
						"Internal error: nesting incorrect; got: %#v",
						currentMessageSpecPart))
				}
				parentGroup.Append(group)
			}
			currentMessageSpecPart = group
			currentNestingLevel = sg.NestingLevel
		} else {
			return nil, errors.New(
				fmt.Sprintf("Parse error at index %d ('%s')",
					index, line))
		}
	}
	return
}

func (p *MessageSpecParser) parseSegmentGroupStart(line string) (segmentGroupStart *SegmentGroupStart, err error) {
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

func (p *MessageSpecParser) parseSegmentEntry(line string) (segmentEntry *SegmentEntry, err error) {
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
func (p *MessageSpecParser) matchHeaderOrEmptyInGroupSection(line string) (matches bool) {
	return strings.HasPrefix(line, "            ")
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
func (p *MessageSpecParser) ParseSpecFile(fileName string) (spec *MessageSpec, err error) {
	// The largest standard message file has 321k (about 6800 lines), so
	// we can read it at once

	log.Printf("Parsing message spec file '%s'", fileName)
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	lines := strings.Split(string(contents), "\n")
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

	specParts, err := p.parseMessageSpecParts(lines[47:])
	if err != nil {
		return
	}

	return NewMessageSpec(id, name, version, release, contrAgency, revision, date, source, specParts), nil
}

func (p *MessageSpecParser) ParseSpecDir(dirName string, suffix string) (specs []*MessageSpec, err error) {
	entries, err := ioutil.ReadDir(dirName)

	specs = []*MessageSpec{}
	for _, entry := range entries {
		fmt.Printf("===== %s", entry.Name())
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

func NewMessageSpecParser(segmentSpecs segment.SegmentSpecProvider) *MessageSpecParser {
	return &MessageSpecParser{
		segmentSpecs: segmentSpecs,
	}
}
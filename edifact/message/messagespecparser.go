package message

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/segment"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const pathSep = string(os.PathSeparator)

var sourceRE = regexp.MustCompile(`^SOURCE: (.*) *$`)

// Parser for message specifications
// e.g. d14b/edmd/AUTHOR_D.14B
type MessageSpecParser struct {
	segmentSpecs segment.SegmentSpecMap
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
	source, err := p.parseSource(lines[46])
	if err != nil {
		return
	}
	return NewMessageSpec(id, name, version, release, contrAgency, revision, date, source), nil
}

func (p *MessageSpecParser) ParseSpecDir(dirName string, suffix string) (specs []*MessageSpec, err error) {
	entries, err := ioutil.ReadDir(dirName)

	specs = []*MessageSpec{}
	for _, entry := range entries {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, "."+suffix) {
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

func NewMessageSpecParser(segmentSpecs segment.SegmentSpecMap) *MessageSpecParser {
	return &MessageSpecParser{
		segmentSpecs: segmentSpecs,
	}
}

package segment

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parses segment specifications file (e.g. EDSD.14B)
type SegmentSpecParser struct {
	compositeElemRE *regexp.Regexp
}

type SegmentSpecMap map[string]*SegmentSpec

// Parse composite element spec of the form
// "020    C138 PRICE MULTIPLIER INFORMATION               C    1"
func (p *SegmentSpecParser) ParseCompositeElemSpec(specStr string) (pos int, name string, title string, mandatory bool, count int, err error) {
	compositeMatch := p.compositeElemRE.FindStringSubmatch(specStr)
	if compositeMatch == nil {
		err = errors.New(fmt.Sprintf("Unable to parse composite spec (specStr: '%s')", specStr))
		return
	}

	if len(compositeMatch) != 6 {
		panic("Internal error: incorrect regular expression")
	}

	pos, err = strconv.Atoi(compositeMatch[1])
	if err != nil {
		return
	}

	name = compositeMatch[2]
	title = strings.TrimSpace(compositeMatch[3])
	mandatoryStr := compositeMatch[4]

	if mandatoryStr == "M" {
		mandatory = true
	} else if mandatoryStr == "C" {
		mandatory = false
	} else {
		err = errors.New("Neither mandatory nor conditional")
		return
	}

	count, err = strconv.Atoi(compositeMatch[5])
	if err != nil {
		return
	}

	return
}

// Parse a single segment specification
func (p *SegmentSpecParser) ParseSpec(specLines []string) (spec *SegmentSpec, err error) {
	/*for _, line := range specLines {

																				}*/
	panic("Not implemented")
}

func (p *SegmentSpecParser) ParseSpecFile(fileName string) (specs SegmentSpecMap, err error) {
	result := SegmentSpecMap{}

	scanner, err := util.NewSpecScanner(fileName)
	if err != nil {
		log.Printf("Unable to create spec scanner for file %s: %s",
			fileName, err)
	}

	first := true

	for {
		// read specification parts
		specLines, err := scanner.GetNextSpecLines(true)

		if err != nil {
			return nil, err
		}

		if !scanner.HasMore && len(specLines) == 0 {
			log.Println("No more lines")
			break
		}

		if first {
			// Skip header part
			first = false
			continue
		}

		spec, err := p.ParseSpec(specLines)
		if err != nil {
			return nil, err
		}
		result[spec.Name] = spec
	}
	return result, nil
}

func NewSegmentSpecParser() *SegmentSpecParser {
	return &SegmentSpecParser{
		compositeElemRE: regexp.MustCompile(`^(\d{3})[ ]{4}(C\d{3}) ([A-Z ]{42}) ([CM])[ ]{4}(\d+)`),
	}
}

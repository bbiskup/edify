package segment

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/dataelement"
	"github.com/bbiskup/edify/edifact/util"
	"regexp"
	"strconv"
	"strings"
)

// Parses segment specifications file (e.g. EDSD.14B)
type SegmentSpecParser struct {
	SimpleDataElemSpecs    dataelement.SimpleDataElementSpecMap
	CompositeDataElemSpecs dataelement.CompositeDataElementSpecMap
	compositeElemRE        *regexp.Regexp
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

	parseSection := func(lines []string) error {
		spec, err := p.ParseSpec(lines)
		if err != nil {
			return err
		}
		result[spec.Id] = spec
		return nil
	}

	err = util.ParseSpecFile(fileName, parseSection)

	return result, err
}

func NewSegmentSpecParser(
	simpleDataElemSpecs dataelement.SimpleDataElementSpecMap,
	compositeDataElemSpecs dataelement.CompositeDataElementSpecMap) *SegmentSpecParser {

	return &SegmentSpecParser{
		SimpleDataElemSpecs:    simpleDataElemSpecs,
		CompositeDataElemSpecs: compositeDataElemSpecs,
		compositeElemRE:        regexp.MustCompile(`^(\d{3})[ ]{4}(C\d{3}) ([A-Z ]{42}) ([CM])[ ]{4}(\d+)`),
	}
}

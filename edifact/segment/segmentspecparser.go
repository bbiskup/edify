package segment

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/dataelement"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type DataElementKind int

const (
	Simple DataElementKind = iota
	Composite
)

// Parses segment specifications file (e.g. EDSD.14B)
type SegmentSpecParser struct {
	SimpleDataElemSpecs    dataelement.SimpleDataElementSpecMap
	CompositeDataElemSpecs dataelement.CompositeDataElementSpecMap
	headerRE               *regexp.Regexp
	dataElemRE             *regexp.Regexp
}

// Parse composite element spec of the form
// "020    C138 PRICE MULTIPLIER INFORMATION               C    1"
func (p *SegmentSpecParser) parseDataElemSpec(
	specStr string) (pos int, id string,
	dataElementKind DataElementKind,
	count int, isMandatory bool,
	err error) {

	dataElemMatch := p.dataElemRE.FindStringSubmatch(specStr)
	if dataElemMatch == nil {
		err = errors.New(fmt.Sprintf("Unable to parse segment spec (specStr: '%s')", specStr))
		return
	}

	if len(dataElemMatch) != 6 {
		panic("Internal error: incorrect regular expression")
	}

	pos, err = strconv.Atoi(dataElemMatch[1])
	if err != nil {
		return
	}

	id = dataElemMatch[2]
	if len(id) == 0 {
		err = errors.New(fmt.Sprintf("Malformed ID: '%s'", id))
	}

	if id[0] == 'C' {
		dataElementKind = Composite
	} else {
		dataElementKind = Simple
	}

	// name = strings.TrimSpace(dataElemMatch[3])
	mandatoryStr := dataElemMatch[4]

	if mandatoryStr == "M" {
		isMandatory = true
	} else if mandatoryStr == "C" {
		isMandatory = false
	} else {
		err = errors.New("Neither mandatory nor conditional")
		return
	}

	count, err = strconv.Atoi(dataElemMatch[5])
	if err != nil {
		return
	}

	return
}

func (p *SegmentSpecParser) parseDataElementSpecs(
	dataElementSpecGroups [][]string) (dataElements []*SegmentDataElementSpec, err error) {

	for _, group := range dataElementSpecGroups {
		if len(group) == 0 {
			return nil, errors.New(fmt.Sprintf("Malformed data element spec group: '%s'", group))
		}
		fmt.Printf("#### group %s\n", group)
		specLine := util.JoinByHangingIndent(group, 8, true)[0]
		if len(specLine) > 0 {
			r, _ := utf8.DecodeRuneInString(specLine)
			if !unicode.IsDigit(r) {
				// May be a note section, e.g.
				//       Note:
				//            The composite C770 - array cell details - occurs
				//            9,999 times in the segment. The use of the ARR

				log.Printf("Skipping group")
				continue
			}
		}

		_, id, kind, count, isMandatory, err := p.parseDataElemSpec(specLine)
		if err != nil {
			return nil, err
		}

		var dataElemSpec dataelement.DataElementSpec
		switch kind {
		case Simple:
			dataElemSpec = p.SimpleDataElemSpecs[id]
		case Composite:
			dataElemSpec = p.CompositeDataElemSpecs[id]
		}

		if dataElemSpec == nil {
			return nil, errors.New(fmt.Sprintf("Data element not found: %s", id))
		}

		segmentDataElementSpec := NewSegmentDataElementSpec(dataElemSpec, count, isMandatory)
		dataElements = append(dataElements, segmentDataElementSpec)
	}
	return
}

func (p *SegmentSpecParser) parseHeader(header string) (id string, name string, err error) {
	headerMatch := p.headerRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New(fmt.Sprintf("Unable to parse header ('%s')", header))
		return
	}

	if len(headerMatch) != 3 {
		panic("Internal error: incorrect regular expression")
	}
	id = headerMatch[1]
	name = headerMatch[2]
	return
}

func (p *SegmentSpecParser) parseFunction(functionLines []string) (fun string, err error) {
	functionPart := util.TrimWhiteSpaceAndJoin(functionLines, " ")
	const functionPrefix = "Function: "
	if strings.HasPrefix(functionPart, functionPrefix) {
		return functionPart[len(functionPrefix):], nil
	} else {
		return "", errors.New(fmt.Sprintf("Unable to parse function: '%s'", functionPart))
	}
}

// Parse a single segment specification
func (p *SegmentSpecParser) ParseSpec(specLines []string) (spec *SegmentSpec, err error) {
	specLines = util.RemoveLeadingAndTrailingEmptyLines(specLines)
	groups := util.SplitMultipleLinesByEmptyLines(specLines)

	headerGroup := groups[0]
	if len(headerGroup) != 1 {
		return nil, errors.New(fmt.Sprintf("Header must have exactly one line: '%s'", headerGroup))
	}
	id, name, err := p.parseHeader(headerGroup[0])
	if err != nil {
		return nil, err
	}
	fun, err := p.parseFunction(groups[1])
	if err != nil {
		return nil, err
	}

	dataElementSpecGroups := groups[2:]
	dataElementSpecs, err := p.parseDataElementSpecs(dataElementSpecGroups)
	if err != nil {
		return
	}

	return NewSegmentSpec(id, name, fun, dataElementSpecs), nil
}

func (p *SegmentSpecParser) ParseSpecFile(fileName string) (specs SegmentSpecProvider, err error) {
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

	return &SegmentSpecProviderImpl{result}, err
}

func NewSegmentSpecParser(
	simpleDataElemSpecs dataelement.SimpleDataElementSpecMap,
	compositeDataElemSpecs dataelement.CompositeDataElementSpecMap) *SegmentSpecParser {

	return &SegmentSpecParser{
		SimpleDataElemSpecs:    simpleDataElemSpecs,
		CompositeDataElemSpecs: compositeDataElemSpecs,
		headerRE:               regexp.MustCompile(`^[ ]{7}([A-Z]{3})  (.*) *$`),
		dataElemRE:             regexp.MustCompile(`^(\d{3})[ ]{4}([0-9C][0-9]{3}) (.+) ([CM])[ ]{3,4}(\d+)[^0-9]?.*`),
	}
}

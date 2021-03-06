package dataelement

import (
	"errors"
	"fmt"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
	"github.com/bbiskup/edify/edifact/spec/specutil"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strings"
)

const (
	ID_IDX = 5
	ID_LEN = 4

	simpleDataElemSectionIndent = 5
)

// Parser for EDIFACT simple data element specification
//
// Sample spec archive:
//    http://www.unece.org/tradewelcome/areas-of-work/un-centre-for-trade-facilitation-and-e-business-uncefact/outputs/standards/unedifact/directories/download.html
// File: EDED.14B
type SimpleDataElemSpecParser struct {
	codesSpecs csp.CodesSpecMap
	numLineRE  *regexp.Regexp
}

func NewSimpleDataElemSpecParser(codesSpecs csp.CodesSpecMap) *SimpleDataElemSpecParser {
	return &SimpleDataElemSpecParser{
		codesSpecs: codesSpecs,
		numLineRE: regexp.MustCompile(
			`^[ ]{4}(\d{4})[ ]+(.*)(\[[BIC]\])$`),
	}
}

// Get data element spec number
func (p *SimpleDataElemSpecParser) getIdAndName(specLinesSections [][]string) (id string, name string, err error) {
	numSection := specLinesSections[0]
	numSectionHeader := numSection[0]
	numLineMatch := p.numLineRE.FindStringSubmatch(numSectionHeader)
	if numLineMatch == nil {
		return "", "", errors.New(
			fmt.Sprintf("Missing num section in line '%s'",
				numSectionHeader))
	}
	id = numLineMatch[1]
	name = strings.TrimSpace(numLineMatch[2])
	return
}

// Get data element description
func (p *SimpleDataElemSpecParser) getDescr(specLinesSections [][]string) (descr string, err error) {
	descLines := specLinesSections[1]
	descLine := util.TrimWhiteSpaceAndJoin(descLines, " ")
	colonIdx := strings.Index(descLine, ":")
	if colonIdx == -1 {
		return "", errors.New("Could not parse description")
	}
	description := strings.TrimSpace(descLine[colonIdx+2:])
	return description, nil
}

func (p *SimpleDataElemSpecParser) getRepr(specLinesSections [][]string) (repr *Repr, err error) {
	reprLine := strings.TrimSpace(specLinesSections[2][0])
	reprLineTokens := strings.Split(reprLine, ":")
	if len(reprLineTokens) != 2 || reprLineTokens[0] != "Repr" {
		return nil, fmt.Errorf("Malformed repr line :'%s'", reprLine)
	}
	reprStr := strings.TrimSpace(reprLineTokens[1])
	return ParseRepr(reprStr)
}

// Parse a single data element spec
func (p *SimpleDataElemSpecParser) ParseSpec(specLines []string) (spec *SimpleDataElemSpec, err error) {
	numSpecLines := len(specLines)
	for i := 0; i < numSpecLines; i++ {
		line := specLines[i]
		if len(line) > 0 {
			specLines[i] = line[1:]
		} else {
			specLines[i] = ""
		}
	}

	specLinesSections := util.SplitByHangingIndent(specLines,
		simpleDataElemSectionIndent-1)
	numSpecLinesSections := len(specLinesSections)
	if numSpecLinesSections < 3 {
		log.Printf("specLines:\n%s\n", strings.Join(specLines, "\n"))
		return nil, fmt.Errorf("Too few (%d) spec segments", numSpecLinesSections)
	}

	id, name, err := p.getIdAndName(specLinesSections)
	if err != nil {
		return nil, err
	}

	description, err := p.getDescr(specLinesSections)
	if err != nil {
		return nil, err
	}

	repr, err := p.getRepr(specLinesSections)
	if err != nil {
		return nil, err
	}

	// may be nil for fields that don't use a code
	codesSpec := p.codesSpecs[id]

	return NewSimpleDataElemSpec(id, name, description, repr, codesSpec)
}

func (p *SimpleDataElemSpecParser) ParseSpecFile(fileName string) (specs SimpleDataElemSpecMap, err error) {
	result := SimpleDataElemSpecMap{}

	parseSection := func(lines []string) error {
		spec, err := p.ParseSpec(lines)
		if err != nil {
			return err
		}
		result[spec.Id()] = spec
		return nil
	}

	err = specutil.ParseSpecFile(fileName, parseSection)
	return result, err
}

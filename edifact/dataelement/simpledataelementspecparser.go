package dataelement

/**
 *  Parser for EDIFACT simple data element specification
 *
 * Sample spec archive:
 *    http://www.unece.org/tradewelcome/areas-of-work/un-centre-for-trade-facilitation-and-e-business-uncefact/outputs/standards/unedifact/directories/download.html
 * File: EDED.14B
 */

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	ID_IDX = 5
	ID_LEN = 4

	simpleDataElementSectionIndent = 5
)

type SimpleDataElementSpecParser struct {
	numLineRE *regexp.Regexp
}

// Get data element spec number
func (p *SimpleDataElementSpecParser) getNameAndNum(specLinesSections [][]string) (name string, num int, err error) {
	numSection := specLinesSections[0]
	numSectionHeader := numSection[0]
	numLineMatch := p.numLineRE.FindStringSubmatch(numSectionHeader)
	if numLineMatch == nil {
		return "", -1, errors.New(
			fmt.Sprintf("Missing num section in line '%s'",
				numSectionHeader))
	}
	name = strings.TrimSpace(numLineMatch[2])
	num, err = strconv.Atoi(numLineMatch[1])
	return
}

// Get data element description
func (p *SimpleDataElementSpecParser) getDescr(specLinesSections [][]string) (descr string, err error) {
	descLine := specLinesSections[1][0]
	colonIdx := strings.Index(descLine, ":")
	if colonIdx == -1 {
		return "", errors.New("Could not parse description")
	}
	description := strings.TrimSpace(descLine[colonIdx:])
	return description, nil
}

func (p *SimpleDataElementSpecParser) getRepr(specLinesSections [][]string) (repr *Repr, err error) {
	reprLine := strings.TrimSpace(specLinesSections[2][0])
	reprLineTokens := strings.Split(reprLine, ":")
	if len(reprLineTokens) != 2 || reprLineTokens[0] != "Repr" {
		return nil, errors.New(fmt.Sprintf("Malformed repr line :'%s'", reprLine))
	}
	reprStr := strings.TrimSpace(reprLineTokens[1])
	return ParseRepr(reprStr)
}

// Parse a single data element spec
func (p *SimpleDataElementSpecParser) ParseSpec(specLines []string) (spec *SimpleDataElementSpec, err error) {
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
		simpleDataElementSectionIndent-1)
	numSpecLinesSections := len(specLinesSections)
	if numSpecLinesSections < 3 {
		log.Printf("specLines:\n%s\n", strings.Join(specLines, "\n"))
		return nil, errors.New(fmt.Sprintf("Too few (%d) spec segments",
			numSpecLinesSections))
	}

	name, num, err := p.getNameAndNum(specLinesSections)
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

	return NewSimpleDataElementSpec(int32(num), name, description, repr), nil
}

type SimpleDataElementSpecMap map[int32]*SimpleDataElementSpec

func (sm SimpleDataElementSpecMap) String() string {
	result := []string{}
	for key, value := range sm {
		result = append(result, fmt.Sprintf("%d: %s", key, value))
	}
	return strings.Join(result, ", ")
}

func (p *SimpleDataElementSpecParser) ParseSpecFile(fileName string) (specs SimpleDataElementSpecMap, err error) {
	result := SimpleDataElementSpecMap{}

	scanner, err := util.NewSpecScanner(fileName)
	if err != nil {
		return
	}

	first := true

	for {
		// read specification parts
		specLines, err := scanner.GetNextSpecLines()

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
		result[spec.Num] = spec
	}
	return result, nil
}

func NewSimpleDataElementSpecParser() *SimpleDataElementSpecParser {
	return &SimpleDataElementSpecParser{
		numLineRE: regexp.MustCompile(
			`^[ ]{4}(\d{4})[ ]+(.*)(\[[BIC]\])$`),
	}
}

package specparser

/**
 *  Parser for EDIFACT specification
 *
 * Sample spec archive:
 *    http://www.unece.org/tradewelcome/areas-of-work/un-centre-for-trade-facilitation-and-e-business-uncefact/outputs/standards/unedifact/directories/download.html
 * File: EDED.14B
 */

import (
	"errors"
	"fmt"
	edi "github.com/bbiskup/edifice/edifact"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	// Separator between specifications (partial)
	specSep = "--------------------"
)

const (
	ID_IDX = 5
	ID_LEN = 4

	dataElementSectionIndent = 5
)

type DataElementSpecParser struct {
	numLineRE *regexp.Regexp
}

// Get data element spec number
func (p *DataElementSpecParser) getNameAndNum(specLinesSections [][]string) (name string, num int, err error) {
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
func (p *DataElementSpecParser) getDescr(specLinesSections [][]string) (descr string, err error) {
	descLine := specLinesSections[1][0]
	colonIdx := strings.Index(descLine, ":")
	if colonIdx == -1 {
		return "", errors.New("Could not parse description")
	}
	description := strings.TrimSpace(descLine[colonIdx:])
	return description, nil
}

func (p *DataElementSpecParser) getRepr(specLinesSections [][]string) (repr *Repr, err error) {
	reprLine := strings.TrimSpace(specLinesSections[2][0])
	reprLineTokens := strings.Split(reprLine, ":")
	if len(reprLineTokens) != 2 || reprLineTokens[0] != "Repr" {
		return nil, errors.New(fmt.Sprintf("Malformed repr line :'%s'", reprLine))
	}
	reprStr := strings.TrimSpace(reprLineTokens[1])
	return ParseRepr(reprStr)
}

// Parse a single data element spec
func (p *DataElementSpecParser) ParseSpec(specLines []string) (spec *DataElementSpec, err error) {
	numSpecLines := len(specLines)
	for i := 0; i < numSpecLines; i++ {
		line := specLines[i]
		if len(line) > 0 {
			specLines[i] = line[1:]
		} else {
			specLines[i] = ""
		}
	}

	specLinesSections := edi.SplitByHangingIndent(specLines,
		dataElementSectionIndent-1)
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

	return NewDataElementSpec(int32(num), name, description, repr), nil
}

type SpecMap map[int32]*DataElementSpec

func (sm SpecMap) String() string {
	result := []string{}
	for key, value := range sm {
		result = append(result, fmt.Sprintf("%d: %s", key, value))
	}
	return strings.Join(result, ", ")
}

func (p *DataElementSpecParser) ParseSpecFile(fileName string) (specs SpecMap, err error) {
	result := SpecMap{}

	scanner, err := NewSpecScanner(fileName)
	if err != nil {
		log.Printf("Unable to create spec scanner for file %s: %s",
			fileName, err)
	}

	first := true

	for {
		// read specification parts
		specLines, err := scanner.GetNextSpecLines()
		// log.Printf("hasMore? %t\n", hasMore)

		if err != nil {
			return nil, err
		}

		if !scanner.hasMore && len(specLines) == 0 {
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

func NewDataElementSpecParser() *DataElementSpecParser {
	return &DataElementSpecParser{
		numLineRE: regexp.MustCompile(
			`^[ ]{4}(\d{4})[ ]+(.*)(\[[BIC]\])$`),
	}
}

package edifact

/**
 *  Parser for EDIFACT specification
 *
 * Sample spec archive:
 *    http://www.unece.org/tradewelcome/areas-of-work/un-centre-for-trade-facilitation-and-e-business-uncefact/outputs/standards/unedifact/directories/download.html
 * File: EDED.14B
 */

import (
	"bufio"
	edi "edifact_experiments/edifact"
	"errors"
	"fmt"
	"log"
	"os"
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

// DataElement specification
type DataElementSpec struct {
	Num  int32
	Name string
	Repr string
}

func (s *DataElementSpec) String() string {
	return fmt.Sprintf("DataElementSpec: %d %s [%s]", s.Num, s.Name, s.Repr)
}

func NewDataElementSpec(num int32, name string, repr string) *DataElementSpec {
	return &DataElementSpec{
		Num:  num,
		Name: name,
		Repr: repr,
	}
}

type DataElementSpecParser struct {
	numLineRE *regexp.Regexp
}

// Parse a single data element spec
func (p *DataElementSpecParser) ParseSpec(specLines []string) (spec *DataElementSpec, err error) {

	specLinesSections := edi.SplitByHangingIndent(specLines, dataElementSectionIndent)
	numSpecLinesSections := len(specLinesSections)
	if numSpecLinesSections < 3 {
		fmt.Printf("specLines:\n%s\n", strings.Join(specLines, "\n"))
		return nil, errors.New(fmt.Sprintf("Too few (%d) spec segments", numSpecLinesSections))
	}

	numSection := specLinesSections[0]
	numSectionHeader := numSection[0]
	numLineMatch := p.numLineRE.FindStringSubmatch(numSectionHeader)
	if numLineMatch == nil {
		return nil, errors.New(fmt.Sprintf("Missing num section in line '%s'", numSectionHeader))
	}
	num, err := strconv.Atoi(numLineMatch[1])
	if err != nil {
		return nil, err
	}

	return NewDataElementSpec(int32(num), "dummyspec", "dummyrepr"), nil
}

// fetch all lines up to next spec separator
func (p *DataElementSpecParser) GetNextSpecLines(scanner *bufio.Scanner) (lines []string, hasMore bool) {
	for {
		scanResult := scanner.Scan()
		if !scanResult {
			if scanner.Err() == nil {
				// EOF
				return lines, false
			}
		}
		err := scanner.Err()
		if err != nil {
			return nil, true
		}

		line := scanner.Text()
		strippedLine := strings.TrimSpace(line)
		if len(strippedLine) == 0 {
			continue
		}

		if strings.HasPrefix(line, specSep) {
			return lines, true
		}

		lines = append(lines, line)
	}
	return lines, true
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
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	first := true

	for {
		// read specification parts
		specLines, hasMore := p.GetNextSpecLines(scanner)
		log.Printf("hasMore? %t\n", hasMore)

		if !hasMore && len(specLines) == 0 {
			log.Println("No more lines")
			break
		}

		if first {
			// Skip header part
			first = false
			continue
		}

		log.Printf("specLines: %s", specLines)

		err := scanner.Err()
		if err != nil {
			return nil, err
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if err != nil {
			return nil, err
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
			`^[ ]{5}(\d{4})[ ]+(.*)(\[[BIC]\])$`),
	}
}

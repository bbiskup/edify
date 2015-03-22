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
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// Separator between specifications (partial)
	specSep = "--------------------"
)

const (
	ID_IDX = 5
	ID_LEN = 4
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
}

// Parse a single data element spec from spec lines
func (p *DataElementSpecParser) ParseSpec(specLines []string) (spec *DataElementSpec, err error) {
	return NewDataElementSpec(100, "dummyspec", "dummyrepr"), nil
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

	for {
		// read specification parts
		specLines, hasMore := p.GetNextSpecLines(scanner)
		log.Printf("hasMore? %t\n", hasMore)

		if !hasMore && len(specLines) == 0 {
			log.Println("No more lines")
			break
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
	return &DataElementSpecParser{}
}
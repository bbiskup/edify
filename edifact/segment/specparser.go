package segment

import (
	"github.com/bbiskup/edifice/edifact/util"
	"log"
)

// Parses segment specifications file (e.g. EDSD.14B)
type SegmentSpecParser struct {
}

type SegmentSpecMap map[string]*SegmentSpec

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
		result[spec.Name] = spec
	}
	return result, nil
}

func NewSegmentSpecParser() *SegmentSpecParser {
	return &SegmentSpecParser{}
}

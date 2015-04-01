package edifact

import (
	"github.com/bbiskup/edify/edifact/codes"
	"github.com/bbiskup/edify/edifact/dataelement"
	"log"
	"os"
	"strings"
)

// Parses all relevant parts of EDIFACT spec
type FullSpecParser struct {
	Version string
}

func NewFullSpecParser(version string, dir string) (*FullSpecParser, error) {
	const sepStr = string(os.PathSeparator)
	codesParser := codes.NewCodesSpecParser()

	codesPath := strings.Join([]string{
		dir, "uncl", "UNCL." + version,
	}, string(os.PathSeparator))

	codesSpecs, err := codesParser.ParseSpecFile(codesPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d codes specs", len(codesSpecs))

	simpleDataElemParser := dataelement.NewSimpleDataElementSpecParser(codesSpecs)

	simpleDataElemSpecPath := strings.Join([]string{
		dir, "eded", "EDED." + version,
	}, string(os.PathSeparator))
	simpleDataElemSpecs, err := simpleDataElemParser.ParseSpecFile(simpleDataElemSpecPath)
	if err != nil {
		return nil, err
	}
	numSimpleDataElemSpecs := len(simpleDataElemSpecs)
	if numSimpleDataElemSpecs > 0 {
		log.Printf("Loaded %d simple data element specs", numSimpleDataElemSpecs)

		// retrieve first element which uses codes (for display)
		var firstVal *dataelement.SimpleDataElementSpec
		for _, v := range simpleDataElemSpecs {
			firstVal = v
			if firstVal.CodesSpecs != nil {
				break
			}
		}
		log.Printf("\tA random spec: %s", firstVal)
		log.Printf("\t  codesSpecs: %d", firstVal.CodesSpecs.Len())
	} else {
		log.Printf("No simple data element specs")
	}

	return &FullSpecParser{version}, nil
}

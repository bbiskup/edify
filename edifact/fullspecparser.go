package edifact

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/codes"
	"github.com/bbiskup/edify/edifact/dataelement"
	"log"
	"os"
	"strings"
)

// Parses all relevant parts of EDIFACT spec
type FullSpecParser struct {
	Version string
	Dir     string
}

func (p *FullSpecParser) getPath(subDir string, filePrefix string) string {
	return strings.Join([]string{
		p.Dir, subDir, fmt.Sprintf("%s.%s", filePrefix, p.Version),
	}, string(os.PathSeparator))
}

func (p *FullSpecParser) parseCodeSpecs() (codes.CodesSpecMap, error) {
	codesParser := codes.NewCodesSpecParser()
	codesPath := p.getPath("uncl", "UNCL")
	codesSpecs, err := codesParser.ParseSpecFile(codesPath)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d codes specs", len(codesSpecs))
	return codesSpecs, nil
}

func (p *FullSpecParser) parseSimpleDataElemSpecs(codesSpecs codes.CodesSpecMap) (dataelement.SimpleDataElementSpecMap, error) {

	simpleDataElemParser := dataelement.NewSimpleDataElementSpecParser(codesSpecs)
	simpleDataElemSpecPath := p.getPath("eded", "EDED")
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
		log.Printf("\tA random spec:")
		log.Printf("%s", firstVal)
		log.Printf("\t  codesSpecs: %d", firstVal.CodesSpecs.Len())
	} else {
		log.Printf("No simple data element specs")
	}
	return simpleDataElemSpecs, nil
}

func (p *FullSpecParser) parseCompositeDataElemSpecs(simpleDataElemSpecs dataelement.SimpleDataElementSpecMap) (dataelement.CompositeDataElementSpecMap, error) {
	compositeDataElemParser := dataelement.NewCompositeDataElementSpecParser(simpleDataElemSpecs)
	compositeDataElemSpecPath := p.getPath("edcd", "EDCD")
	compositeDataElemSpecs, err := compositeDataElemParser.ParseSpecFile(compositeDataElemSpecPath)
	if err != nil {
		return nil, err
	}

	numCompositeDataElemSpecs := len(compositeDataElemSpecs)
	if numCompositeDataElemSpecs > 0 {
		log.Printf("Loaded %d composite data element specs", numCompositeDataElemSpecs)
	}
	return compositeDataElemSpecs, nil
}

func (p *FullSpecParser) Parse() error {
	codeSpecs, err := p.parseCodeSpecs()
	if err != nil {
		return err
	}

	simpleDataElemSpecs, err := p.parseSimpleDataElemSpecs(codeSpecs)
	if err != nil {
		return err
	}

	compositeDataElemSpecs, err := p.parseCompositeDataElemSpecs(simpleDataElemSpecs)
	if err != nil {
		return err
	}

	_ = compositeDataElemSpecs
	return nil
}

func NewFullSpecParser(version string, dir string) (*FullSpecParser, error) {
	return &FullSpecParser{version, dir}, nil
}

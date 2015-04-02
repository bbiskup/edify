package edifact

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/codes"
	"github.com/bbiskup/edify/edifact/dataelement"
	"github.com/bbiskup/edify/edifact/segment"
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
	parser := codes.NewCodesSpecParser()
	path := p.getPath("uncl", "UNCL")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d codes specs", len(specs))
	return specs, nil
}

func (p *FullSpecParser) parseSimpleDataElemSpecs(codesSpecs codes.CodesSpecMap) (dataelement.SimpleDataElementSpecMap, error) {

	parser := dataelement.NewSimpleDataElementSpecParser(codesSpecs)
	path := p.getPath("eded", "EDED")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}
	numSpecs := len(specs)
	if numSpecs > 0 {
		log.Printf("Loaded %d simple data element specs", numSpecs)

		// retrieve first element which uses codes (for display)
		var firstVal *dataelement.SimpleDataElementSpec
		for _, v := range specs {
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
	return specs, nil
}

func (p *FullSpecParser) parseCompositeDataElemSpecs(simpleDataElemSpecs dataelement.SimpleDataElementSpecMap) (dataelement.CompositeDataElementSpecMap, error) {
	parser := dataelement.NewCompositeDataElementSpecParser(simpleDataElemSpecs)
	path := p.getPath("edcd", "EDCD")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}

	numSpecs := len(specs)
	if numSpecs > 0 {
		log.Printf("Loaded %d composite data element specs", numSpecs)
	}
	return specs, nil
}

func (p *FullSpecParser) parseSegmentSpecs(
	simpleDataElemSpecs dataelement.SimpleDataElementSpecMap,
	compositeDataElemSpecs dataelement.CompositeDataElementSpecMap) (segment.SegmentSpecMap, error) {

	parser := segment.NewSegmentSpecParser(simpleDataElemSpecs, compositeDataElemSpecs)
	path := p.getPath("edsd", "EDSD")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}

	numSpecs := len(specs)
	if numSpecs > 0 {
		log.Printf("Loaded %d segment specs", numSpecs)
	}
	return specs, nil
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

	segmentSpecs, err := p.parseSegmentSpecs(simpleDataElemSpecs, compositeDataElemSpecs)

	_ = segmentSpecs
	return err
}

func NewFullSpecParser(version string, dir string) (*FullSpecParser, error) {
	return &FullSpecParser{version, dir}, nil
}

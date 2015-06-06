package specparser

import (
	"errors"
	"fmt"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
	"os"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

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

func (p *FullSpecParser) parseCodeSpecs() (csp.CodesSpecMap, error) {
	parser := csp.NewCodesSpecParser()
	path := p.getPath("uncl", "UNCL")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("Loaded %d codes specs", len(specs))
	return specs, nil
}

func (p *FullSpecParser) parseSimpleDataElemSpecs(codesSpecs csp.CodesSpecMap) (dsp.SimpleDataElemSpecMap, error) {

	parser := dsp.NewSimpleDataElemSpecParser(codesSpecs)
	path := p.getPath("eded", "EDED")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}
	numSpecs := len(specs)
	if numSpecs > 0 {
		log.Printf("Loaded %d simple data element specs", numSpecs)

		// retrieve first element which uses codes (for display)
		var firstVal *dsp.SimpleDataElemSpec
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

func (p *FullSpecParser) parseCompositeDataElemSpecs(simpleDataElemSpecs dsp.SimpleDataElemSpecMap) (dsp.CompositeDataElemSpecMap, error) {
	parser := dsp.NewCompositeDataElemSpecParser(simpleDataElemSpecs)
	path := p.getPath("edcd", "EDCD")
	specs, err := parser.ParseSpecFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error parsing composite data elem file %s: %s", path, err))
	}

	numSpecs := len(specs)
	if numSpecs > 0 {
		log.Printf("Loaded %d composite data element specs", numSpecs)
	}
	return specs, nil
}

func (p *FullSpecParser) parseSegSpecs(
	simpleDataElemSpecs dsp.SimpleDataElemSpecMap,
	compositeDataElemSpecs dsp.CompositeDataElemSpecMap) (specs ssp.SegSpecProvider, err error) {

	parser := ssp.NewSegSpecParser(simpleDataElemSpecs, compositeDataElemSpecs)
	path := p.getPath("edsd", "EDSD")
	specs, err = parser.ParseSpecFile(path)
	if err != nil {
		return nil, err
	}

	numSpecs := specs.Len()
	if numSpecs > 0 {
		log.Printf("Loaded %d segment specs", numSpecs)
	}
	return specs, nil
}

func (p *FullSpecParser) ParseMsgSpecs(segSpecs ssp.SegSpecProvider) (msgSpecs msp.MsgSpecMap, err error) {
	msgDir := p.Dir + pathSeparator + "edmd"
	parser := msp.NewMsgSpecParser(segSpecs)
	fmt.Printf("Parsing message specs with suffix '%s' in directory '%s'", p.Version, msgDir)
	return parser.ParseSpecDir(msgDir, p.Version)
}

func (p *FullSpecParser) ParseSegSpecsWithPrerequisites() (ssp.SegSpecProvider, error) {
	codeSpecs, err := p.parseCodeSpecs()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error parsing code specs: %s", err))
	}

	simpleDataElemSpecs, err := p.parseSimpleDataElemSpecs(codeSpecs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error parsing simple data element specs: %s", err))
	}

	compositeDataElemSpecs, err := p.parseCompositeDataElemSpecs(simpleDataElemSpecs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Error parsing composite data element specs: %s", err))
	}

	return p.parseSegSpecs(simpleDataElemSpecs, compositeDataElemSpecs)
}

func (p *FullSpecParser) Parse() error {

	segSpecs, err := p.ParseSegSpecsWithPrerequisites()
	if err != nil {
		return err
	}

	msgSpecs, err := p.ParseMsgSpecs(segSpecs)
	if err != nil {
		return err
	}

	log.Printf("Parsed %d message specs", len(msgSpecs))
	return err
}

func NewFullSpecParser(version string, dir string) (*FullSpecParser, error) {
	return &FullSpecParser{version, dir}, nil
}

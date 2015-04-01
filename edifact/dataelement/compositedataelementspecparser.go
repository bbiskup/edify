package dataelement

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parses composite element specification (e.g. EDSD.14B)
type CompositeDataElementSpecParser struct {
	headerRE        *regexp.Regexp
	componentElemRE *regexp.Regexp
}

/* TODO: move to segment parser
// Parse composite element spec header of the form
// "020    C138 PRICE MULTIPLIER INFORMATION               C    1"
func (p *CompositeDataElementSpecParser) ParseHeader(header string) (pos int, name string, title string, mandatory bool, count int, err error) {
	headerMatch := p.headerRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New(fmt.Sprintf("Missing header section (header: '%s'", header))
		return
	}

	if len(headerMatch) != 6 {
		panic("Internal error: incorrect regular expression")
	}

	pos, err = strconv.Atoi(headerMatch[1])
	if err != nil {
		return
	}

	name = headerMatch[2]
	title = strings.TrimSpace(headerMatch[3])
	mandatoryStr := headerMatch[4]

	if mandatoryStr == "M" {
		mandatory = true
	} else if mandatoryStr == "C" {
		mandatory = false
	} else {
		err = errors.New("Neither mandatory nor conditional")
		return
	}

	count, err = strconv.Atoi(headerMatch[5])
	if err != nil {
		return
	}

	return
}*/

// Parse a line of the form
// "       3477  Address format code                       M      an..3"
func (p *CompositeDataElementSpecParser) ParseComponentDataElemenSpec(specLine string) (spec *ComponentDataElementSpec, err error) {
	specMatch := p.componentElemRE.FindStringSubmatch(specLine)
	if specMatch == nil {
		err = errors.New(fmt.Sprintf("Failed to match component data element spec '%s'", specLine))
		return
	}

	if len(specMatch) != 4 {
		panic("Internal error: incorrect regular expression")
	}

	numStr := specMatch[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return
	}

	mandatoryStr := specMatch[3]
	var isMandatory bool

	if mandatoryStr == "M" {
		isMandatory = true
	} else if mandatoryStr == "C" {
		isMandatory = false
	} else {
		err = errors.New("Neither mandatory nor conditional")
		return
	}

	spec = NewComponentDataElementSpec(int32(num), isMandatory)
	return
}

/* Parse spec for a single composite data element of the form

020    C138 PRICE MULTIPLIER INFORMATION               C    1
       5394  Price multiplier rate                     M      n..12
       5393  Price multiplier type code qualifier      C      an..3

*/ /*
func (p *CompositeDataElementSpecParser) ParseSpec(specLines []string) (spec *CompositeDataElementSpec, err error) {
	if len(specLines) == 0 {
		return nil, errors.New("Missing composite spec header")
	}

	_, name, title, isMandatory, count, err := p.ParseHeader(specLines[0])
	if err != nil {
		return
	}

	componentSpecs := []*ComponentDataElementSpec{}
	numSpecLines := len(specLines)
	for i := 1; i < numSpecLines; i++ {
		line := specLines[i]
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		elemSpec, err := p.ParseComponentDataElemenSpec(line)
		if err != nil {
			return nil, err
		}
		componentSpecs = append(componentSpecs, elemSpec)
	}

	if len(componentSpecs) == 0 {
		return nil, errors.New("No component elements")
	}

	spec = NewCompositeDataElementSpec(name, title, count, isMandatory, componentSpecs)
	return
}*/

func (p *CompositeDataElementSpecParser) ParseHeader(header string) (id string, name string, err error) {
	headerMatch := p.headerRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New(fmt.Sprintf("Missing header section (header: '%s'", header))
		return
	}

	if len(headerMatch) != 3 {
		panic("Internal error: incorrect regular expression")
	}

	id = headerMatch[1]
	name = headerMatch[2]
	return
}

/** Parse multiple component specification lines of the form

010    5105  Monetary amount function detail
           description code                          C      an..17
020    1131  Code list identification code             C      an..17
030    3055  Code list responsible agency code         C      an..3
040    5104  Monetary amount function detail
           description                               C      an..70
*/
func (p *CompositeDataElementSpecParser) ParseComponentSpecs(componentGroup []string) ([]*ComponentDataElementSpec, error) {
	componentSpecs := []*ComponentDataElementSpec{}

	joinedSpecs := util.JoinByHangingIndent(componentGroup, 0, true)
	log.Printf("Joined specs: %#v", joinedSpecs)

	numComponents := len(joinedSpecs)

	if numComponents == 0 {
		return nil, errors.New("No component elements")
	}

	for i := 1; i < numComponents; i++ {
		line := joinedSpecs[i]
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		elemSpec, err := p.ParseComponentDataElemenSpec(line)
		if err != nil {
			return nil, err
		}
		componentSpecs = append(componentSpecs, elemSpec)
	}

	log.Printf("Found %d component specs", len(componentSpecs))
	return componentSpecs, nil
}

/**
  Parse a single section from composite data element spec file, e.g. EDCD14B
  Example:

         C001 TRANSPORT MEANS

         Desc: Code and/or name identifying the type of means of
               transport.

  010    8179  Transport means description code          C      an..8
  020    1131  Code list identification code             C      an..17
  030    3055  Code list responsible agency code         C      an..3
  040    8178  Transport means description               C      an..17
*/
func (p *CompositeDataElementSpecParser) ParseSpec(specLines []string) (spec *CompositeDataElementSpec, err error) {
	if len(specLines) == 0 {
		return nil, errors.New("Missing composite spec header")
	}

	log.Printf("specLines: \n%s\n", specLines)
	specLines = util.RemoveLeadingAndTrailingEmptyLines(specLines)
	groups := util.SplitMultipleLinesByEmptyLines(specLines)
	log.Printf("Groups: \n%s\n", groups)

	if len(groups) < 3 {
		return nil, errors.New(fmt.Sprintf("Not enough groups for spec %s", groups))
	}

	headerGroup := groups[0]
	if len(headerGroup) != 1 {
		return nil, errors.New(fmt.Sprintf("Header group must contain a single line (%s)", headerGroup))
	}
	header := headerGroup[0]

	id, name, err := p.ParseHeader(header)
	if err != nil {
		return
	}

	// TODO parse description
	_ = groups[1]

	componentGroup := groups[2]

	componentSpecs, err := p.ParseComponentSpecs(componentGroup)
	if err != nil {
		return
	}

	spec = NewCompositeDataElementSpec(id, name, "dummy description", componentSpecs)

	return
}

type CompositeDataElementSpecMap map[string]*CompositeDataElementSpec

func (sm CompositeDataElementSpecMap) String() string {
	result := []string{}
	for key, value := range sm {
		result = append(result, fmt.Sprintf("%s: %s", key, value))
	}
	return strings.Join(result, ", ")
}

// Parse composite data element spec file, e.g. EDCD14.B
func (p *CompositeDataElementSpecParser) ParseSpecFile(fileName string) (specs CompositeDataElementSpecMap, err error) {
	result := CompositeDataElementSpecMap{}

	scanner, err := util.NewSpecScanner(fileName)
	if err != nil {
		return
	}

	first := true

	for {
		// read specification parts
		specLines, err := scanner.GetNextSpecLines(false)

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

		log.Printf("#### specLines: \n%s\n", specLines)
		spec, err := p.ParseSpec(specLines)
		if err != nil {
			return nil, err
		}
		result[spec.Name] = spec
	}
	return result, nil
}

func NewCompositeDataElementSpecParser() *CompositeDataElementSpecParser {
	// TODO account for change indicators?
	return &CompositeDataElementSpecParser{
		// We ignore the repr spec (already available via simple data element specs)
		// TODO: move header RE to segment parser
		//headerRE:        regexp.MustCompile(`^(\d{3})[ ]{4}(C\d{3}) ([A-Z ]{42}) ([CM])[ ]{4}(\d+)`),

		headerRE: regexp.MustCompile(`^[ ]{7}(C[0-9]{3}) ([A-Z/& -]+) *$`),
		// OBSOLETE componentElemRE: regexp.MustCompile(`^[ ]{7}(\d{4})  ([A-Za-z ]{41}) ([CM]) .+$`),
		componentElemRE: regexp.MustCompile(`^([0-9]{3})[ ]{4}([0-9]{4})  [A-Za-z- ]+ ([CM]) .+$`),
	}
}

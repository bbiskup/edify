package dataelement

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parses composite element specification (e.g. EDSD.14B)
type CompositeDataElementSpecParser struct {
	headerRE        *regexp.Regexp
	componentElemRE *regexp.Regexp
}

// Parse composite element spec header of the form
// "020    C138 PRICE MULTIPLIER INFORMATION               C    1"
func (p *CompositeDataElementSpecParser) ParseHeader(header string) (pos int, name string, title string, mandatory bool, count int, err error) {
	headerMatch := p.headerRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New("Missing header section")
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
}

// Parse a line of the form
// "       3477  Address format code                       M      an..3"
func (p *CompositeDataElementSpecParser) ParseComponentDataElementSpec(specLine string) (spec *ComponentDataElementSpec, err error) {
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

*/
func (p *CompositeDataElementSpecParser) Parse(specLines []string) (spec *CompositeDataElementSpec, err error) {
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
		elemSpec, err := p.ParseComponentDataElementSpec(line)
		if err != nil {
			return nil, err
		}
		componentSpecs = append(componentSpecs, elemSpec)
	}

	spec = NewCompositeDataElementSpec(name, title, count, isMandatory, componentSpecs)
	return
}

func NewCompositeDataElementSpecParser() *CompositeDataElementSpecParser {
	// TODO account for change indicators?
	return &CompositeDataElementSpecParser{
		// We ignore the repr spec (already available via simple data element specs)
		headerRE:        regexp.MustCompile(`^(\d{3})[ ]{4}(C\d{3}) ([A-Z ]{42}) ([CM])[ ]{4}(\d+)`),
		componentElemRE: regexp.MustCompile(`^[ ]{7}(\d{4})  ([A-Za-z ]{41}) ([CM]) .+$`),
	}
}

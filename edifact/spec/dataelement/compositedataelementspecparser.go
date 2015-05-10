package dataelement

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/specutil"
	"github.com/bbiskup/edify/edifact/util"
	"regexp"
	"strconv"
	"strings"
)

// Parses composite element specification (e.g. EDSD.14B)
type CompositeDataElementSpecParser struct {
	simpleDataElemSpecs SimpleDataElementSpecMap
	headerRE            *regexp.Regexp
	componentElemRE     *regexp.Regexp
}

/* Parse a line of the form
010    3148  Communication address identifier          M      an..512
*/
func (p *CompositeDataElementSpecParser) ParseComponentDataElemenSpec(specLine string) (spec *ComponentDataElementSpec, err error) {
	specMatch := p.componentElemRE.FindStringSubmatch(specLine)
	if specMatch == nil {
		err = errors.New(
			fmt.Sprintf(
				"Failed to match component data element spec '%s'",
				specLine))
		return
	}

	if len(specMatch) != 4 {
		panic("Internal error: incorrect regular expression")
	}

	positionStr := specMatch[1]
	position, err := strconv.Atoi(positionStr)
	if err != nil {
		return nil, err
	}

	id := specMatch[2]

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

	simpleDataElemSpec := p.simpleDataElemSpecs[id]
	if simpleDataElemSpec == nil {
		return nil, errors.New(fmt.Sprintf("No simple data elem spec for ID %s", id))
	}
	spec = NewComponentDataElementSpec(position, isMandatory, simpleDataElemSpec)
	return
}

func (p *CompositeDataElementSpecParser) ParseHeader(header string) (id string, name string, err error) {
	headerMatch := p.headerRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New(
			fmt.Sprintf(
				"Missing header section (header: '%s'",
				header))
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

	// log.Printf("Found %d component specs", len(componentSpecs))
	return componentSpecs, nil
}

/**
Parse a composite data element description of the form

       Desc: To specify the event category.

or multi-line:

       Desc: Communication number of a department or employee in
             a specified channel.

*/
func (p *CompositeDataElementSpecParser) ParseDescription(lines []string) (string, error) {
	lines = util.JoinByHangingIndent(lines, 7, true)
	if len(lines) != 1 {
		return "", errors.New(
			fmt.Sprintf(
				"Failed to parse description '%s'",
				lines))
	}
	line := lines[0]

	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "Desc: ") {
		line = line[6:]
	}
	return line, nil
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

	// log.Printf("specLines: \n%s\n", specLines)
	specLines = util.RemoveLeadingAndTrailingEmptyLines(specLines)
	groups := util.SplitMultipleLinesByEmptyLines(specLines)
	// log.Printf("Groups: \n%s\n", groups)

	if len(groups) < 3 {
		return nil, errors.New(fmt.Sprintf("Not enough groups for spec %s", groups))
	}

	headerGroup := groups[0]
	if len(headerGroup) != 1 {
		return nil, errors.New(
			fmt.Sprintf(
				"Header group must contain a single line (%s)",
				headerGroup))
	}
	header := headerGroup[0]

	id, name, err := p.ParseHeader(header)
	if err != nil {
		return
	}

	descrGroup := groups[1]
	var descr string
	if len(descrGroup) == 0 {
		descr = "<no description>"
	} else {
		descr, err = p.ParseDescription(descrGroup)
		if err != nil {
			return nil, err
		}
	}

	componentGroup := groups[2]

	componentSpecs, err := p.ParseComponentSpecs(componentGroup)
	if err != nil {
		return
	}

	spec = NewCompositeDataElementSpec(id, name, descr, componentSpecs)

	return
}

// Parse composite data element spec file, e.g. EDCD14.B
func (p *CompositeDataElementSpecParser) ParseSpecFile(fileName string) (specs CompositeDataElementSpecMap, err error) {
	result := CompositeDataElementSpecMap{}

	parseSection := func(lines []string) error {
		spec, err := p.ParseSpec(lines)
		if err != nil {
			return err
		}
		result[spec.Id()] = spec
		return nil
	}

	err = specutil.ParseSpecFile(fileName, parseSection)

	return result, err
}

func NewCompositeDataElementSpecParser(simpleDataElemSpecs SimpleDataElementSpecMap) *CompositeDataElementSpecParser {
	// TODO account for change indicators?
	return &CompositeDataElementSpecParser{
		// We ignore the repr spec (already available via simple data element specs)

		simpleDataElemSpecs: simpleDataElemSpecs,

		headerRE: regexp.MustCompile(`^[ ]{7}(C[0-9]{3}) ([A-Z/& -]+) *$`),
		// OBSOLETE componentElemRE: regexp.MustCompile(`^[ ]{7}(\d{4})  ([A-Za-z ]{41}) ([CM]) .+$`),
		componentElemRE: regexp.MustCompile(
			`^([0-9]{3})[ ]{4}([0-9]{4})  [A-Za-z- ]+ ([CM]) .+$`),
	}
}

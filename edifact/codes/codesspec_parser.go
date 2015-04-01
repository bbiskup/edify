package codes

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parser for EDIFACT code list (e.g. UNCL.14B)
type CodesSpecParser struct {
	codeHeaderRE  *regexp.Regexp
	codesHeaderRE *regexp.Regexp
}

/**
  Parse a header of the form:

     1073  Document line action code                               [B]

*/
func (p *CodesSpecParser) ParseCodesSpecHeader(header string) (id int32, name string, err error) {
	headerMatch := p.codesHeaderRE.FindStringSubmatch(header)
	if headerMatch == nil {
		err = errors.New(fmt.Sprintf("Unable to parse codes header section (header: '%s'", header))
		return
	}

	if len(headerMatch) != 4 {
		panic("Internal error: incorrect regular expression")
	}

	id_, err := strconv.Atoi(headerMatch[1])
	if err != nil {
		return
	}
	id = int32(id_)
	name = headerMatch[2]
	return
}

/**
Parse a description of the form

     Desc: Code specifying a section of a message.

or multi-line:

     Desc: Code indicating an action associated with a line of a
           document.

*/
func (p *CodesSpecParser) ParseDescription(lines []string) (string, error) {
	line := util.TrimWhiteSpaceAndJoin(lines, " ")

	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "Desc: ") {
		line = line[6:]
	}
	return line, nil
}

func (p *CodesSpecParser) ParseCodeSpecHeader(header string) (id string, name string, err error) {
	codeHeaderMatch := p.codeHeaderRE.FindStringSubmatch(header)
	if codeHeaderMatch == nil {
		err = errors.New(fmt.Sprintf("Unable to parse code header section (header: '%s'", header))
		return "", "", err
	}

	if len(codeHeaderMatch) != 3 {
		panic("Internal error: incorrect regular expression")
	}
	id = codeHeaderMatch[1]
	name = codeHeaderMatch[2]
	return
}

/**
Parse single code spec of the form

     7     Sub-line item
              The section of the message being referenced refers to
              the sub-line item.
*/
func (p *CodesSpecParser) ParseCodeSpec(specLines []string) (*CodeSpec, error) {
	if len(specLines) < 2 {
		return nil, errors.New("Missing spec header and/or description")
	}
	id, name, err := p.ParseCodeSpecHeader(specLines[0])
	if err != nil {
		return nil, err
	}

	descriptionLines := specLines[1:]
	description := util.TrimWhiteSpaceAndJoin(descriptionLines, " ")

	return NewCodeSpec(id, name, description), nil
}

/* Parse multiple code specifications of the form

   7     Structure hierarchical link
            The list contains a set of hierarchical link values
            between structures.

   8     Structure group link
            The list contains a set of group link values between
            structures.

   9     Multiple hierarchical structure item
            The list contains a set of items at multiple
            hierarchical levels in a structure.

*/
func (p *CodesSpecParser) ParseCodeSpecs(lines []string) ([]*CodeSpec, error) {
	result := []*CodeSpec{}
	groups := util.SplitByHangingIndent(lines, 5)

	for _, group := range groups {
		spec, err := p.ParseCodeSpec(group)
		if err != nil {
			return nil, err
		}
		result = append(result, spec)
	}
	return result, nil
}

func (p *CodesSpecParser) ParseCodesSpec(specLines []string) (spec *CodesSpec, err error) {
	groups := util.SplitMultipleLinesByEmptyLines(specLines[1:])
	log.Printf("Groups: \n%s\n", groups)

	if len(groups) < 4 {
		return nil, errors.New(fmt.Sprintf("Not enough groups for spec %s", groups))
	}

	headerGroup := groups[0]
	if len(headerGroup) == 0 {
		return nil, errors.New(fmt.Sprintf("Missing header (%s)", headerGroup))
	}
	if headerGroup[0][0] != ' ' {
		// Change mark
		headerGroup[0] = headerGroup[0][1:]
	}
	header := util.TrimWhiteSpaceAndJoin(headerGroup, " ")

	id, name, err := p.ParseCodesSpecHeader(header)
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

	codeSpecs, err := p.ParseCodeSpecs(groups[3])
	if err != nil {
		return
	}

	return NewCodesSpec(id, name, descr, codeSpecs), nil
}

type CodesSpecMap map[int32]*CodesSpec

func (p *CodesSpecParser) ParseSpecFile(fileName string) (specs CodesSpecMap, err error) {
	result := CodesSpecMap{}

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
		spec, err := p.ParseCodesSpec(specLines)
		if err != nil {
			return nil, err
		}
		result[spec.Id] = spec
	}
	return result, nil
}

func NewCodesSpecParser() *CodesSpecParser {
	return &CodesSpecParser{
		codeHeaderRE:  regexp.MustCompile(`^[ ]{5}(.{1,3})[ ]+(.*) *$`),
		codesHeaderRE: regexp.MustCompile(`^([0-9]{4})  (.+) (\[[BIC]\]) *$`),
	}
}

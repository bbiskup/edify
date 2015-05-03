package codes

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"strings"
)

// EDIFACT Codes as defined e.g. in UNCL.14B
type CodesSpec struct {
	Id          string
	Name        string
	Description string
	CodeSpecs   []*CodeSpec
}

type CodesSpecMap map[string]*CodesSpec

func (s *CodesSpec) String() string {
	specsStrs := []string{}
	for _, spec := range s.CodeSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	codeSpecsStr := strings.Join(specsStrs, "\n")
	descriptionStr := util.Ellipsis(s.Description, maxDescrDisplayLen)

	return fmt.Sprintf("%s %s %s\n%s", s.Id, s.Name, descriptionStr, codeSpecsStr)
}

func (s *CodesSpec) Len() int {
	if s.CodeSpecs == nil {
		return 0
	} else {
		return len(s.CodeSpecs)
	}
}

func NewCodesSpec(id string, name string, description string, codeSpecs []*CodeSpec) *CodesSpec {
	return &CodesSpec{
		Id:          id,
		Name:        name,
		Description: description,
		CodeSpecs:   codeSpecs,
	}
}

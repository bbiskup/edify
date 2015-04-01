package codes

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"strings"
)

// EDIFACT Codes as defined e.g. in UNCL.14B
type CodesSpec struct {
	Id          int32
	Name        string
	Description string
	CodeSpecs   []*CodeSpec
}

func (s *CodesSpec) String() string {
	specsStrs := []string{}
	for _, spec := range s.CodeSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	codeSpecsStr := strings.Join(specsStrs, "\n")

	descriptionStr := util.Ellipsis(s.Description, maxDescrDisplayLen)

	return fmt.Sprintf("%d %s %s\n%s", s.Id, s.Name, descriptionStr, codeSpecsStr)
}

func NewCodesSpec(id int32, name string, codeSpecs []*CodeSpec) *CodesSpec {
	return &CodesSpec{
		Id:        id,
		Name:      name,
		CodeSpecs: codeSpecs,
	}
}

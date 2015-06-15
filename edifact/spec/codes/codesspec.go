package codes

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"sort"
	"strings"
)

type CodeSpecMap map[string]*CodeSpec

// EDIFACT Codes as defined e.g. in UNCL.14B
type CodesSpec struct {
	Id          string
	Name        string
	Description string
	codeSpecMap CodeSpecMap
	CodeSpecs   []*CodeSpec
}

func NewCodesSpec(id string, name string, description string, codeSpecs []*CodeSpec) *CodesSpec {
	codeSpecMap := CodeSpecMap{}
	for _, codeSpec := range codeSpecs {
		codeSpecMap[codeSpec.Id] = codeSpec
	}

	return &CodesSpec{
		Id:          id,
		Name:        name,
		Description: description,
		codeSpecMap: codeSpecMap,
		CodeSpecs:   codeSpecs,
	}
}

type CodesSpecMap map[string]*CodesSpec

func (s *CodesSpec) Contains(code string) bool {
	return s.codeSpecMap[code] != nil
}

func (s *CodesSpec) String() string {
	specsStrs := []string{}
	for _, spec := range s.codeSpecMap {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	codeSpecsStr := strings.Join(specsStrs, "\n")
	descriptionStr := util.Ellipsis(s.Description, maxDescrDisplayLen)

	return fmt.Sprintf("%s %s %s\n%s", s.Id, s.Name, descriptionStr, codeSpecsStr)
}

func (s *CodesSpec) CodeListStr() string {
	codes := make(sort.StringSlice, 0, len(s.codeSpecMap))
	for code, _ := range s.codeSpecMap {
		codes = append(codes, code)
	}
	codes.Sort()
	return strings.Join(codes, ", ")
}

func (s *CodesSpec) Len() int {
	if s.codeSpecMap == nil {
		return 0
	} else {
		return len(s.codeSpecMap)
	}
}

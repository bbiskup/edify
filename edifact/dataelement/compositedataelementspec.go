package dataelement

import (
	"fmt"
	"strings"
)

type CompositeDataElementSpec struct {
	Num            string
	Name           string
	ComponentSpecs []*SimpleDataElementSpec
}

func (s *CompositeDataElementSpec) String() string {
	specsStrs := []string{}
	for _, spec := range s.ComponentSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec))
	}
	componentSpecsStr := strings.Join(specsStrs, "\n")
	return fmt.Sprintf("%s %s\n%s", s.Num, s.Name, componentSpecsStr)
}

func NewCompositeDataElementSpec(num string, name string, componentSpecs []*SimpleDataElementSpec) *CompositeDataElementSpec {
	return &CompositeDataElementSpec{
		Num:            num,
		Name:           name,
		ComponentSpecs: componentSpecs,
	}
}

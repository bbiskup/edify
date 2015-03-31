package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"strings"
)

type CompositeDataElementSpec struct {
	Name           string
	Title          string
	Count          int
	IsMandatory    bool
	ComponentSpecs []*ComponentDataElementSpec
}

func (s *CompositeDataElementSpec) String() string {
	specsStrs := []string{}
	for _, spec := range s.ComponentSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	componentSpecsStr := strings.Join(specsStrs, "\n")
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mandatory", "conditional")
	return fmt.Sprintf(
		"Composite %s %s %d (%s)\n%s",
		s.Name, s.Title, s.Count, isMandatoryStr, componentSpecsStr)
}

func NewCompositeDataElementSpec(
	name string, title string, count int, isMandatory bool,
	componentSpecs []*ComponentDataElementSpec) *CompositeDataElementSpec {

	return &CompositeDataElementSpec{
		Name:           name,
		Title:          title,
		Count:          count,
		IsMandatory:    isMandatory,
		ComponentSpecs: componentSpecs,
	}
}

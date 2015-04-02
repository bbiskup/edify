package dataelement

import (
	"bytes"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"strings"
)

const (
	maxDescrDisplayLen = 10
)

type CompositeDataElementSpec struct {
	id             string
	name           string
	Description    string
	ComponentSpecs []*ComponentDataElementSpec
}

// from interface DataElementSpec
func (s *CompositeDataElementSpec) Id() string {
	return s.id
}

// from interface DataElementSpec
func (s *CompositeDataElementSpec) Name() string {
	return s.name
}

func (s *CompositeDataElementSpec) String() string {
	descrStr := util.Ellipsis(s.Description, maxDescrDisplayLen)
	specsStrs := []string{}
	for _, spec := range s.ComponentSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	componentSpecsStr := strings.Join(specsStrs, "\n")
	return fmt.Sprintf(
		"Composite %s %s '%s'\n%s",
		s.id, s.name, descrStr, componentSpecsStr)
}

func NewCompositeDataElementSpec(
	id string, name string, description string,
	componentSpecs []*ComponentDataElementSpec) *CompositeDataElementSpec {

	return &CompositeDataElementSpec{
		id:             id,
		name:           name,
		Description:    description,
		ComponentSpecs: componentSpecs,
	}
}

type CompositeDataElementSpecMap map[string]*CompositeDataElementSpec

/*
func (sm CompositeDataElementSpecMap) String() string {
	result := []string{}
	for key, value := range sm {
		result = append(result, fmt.Sprintf("%s: %s", key, value))
	}
	return strings.Join(result, ", ")
}*/

func (m CompositeDataElementSpecMap) String() string {
	var result bytes.Buffer
	result.WriteString("CompositeDataElementSpecMap\n")
	for id, spec := range m {
		result.WriteString(fmt.Sprintf("\t%-8s: %s\n", id, spec))
	}
	return result.String()
}

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

type CompositeDataElemSpec struct {
	id             string
	name           string
	Descr          string
	ComponentSpecs []*ComponentDataElemSpec
}

func NewCompositeDataElemSpec(
	id string, name string, descr string,
	componentSpecs []*ComponentDataElemSpec) *CompositeDataElemSpec {

	return &CompositeDataElemSpec{
		id:             id,
		name:           name,
		Descr:          descr,
		ComponentSpecs: componentSpecs,
	}
}

func (m CompositeDataElemSpecMap) String() string {
	var result bytes.Buffer
	result.WriteString("CompositeDataElemSpecMap\n")
	for id, spec := range m {
		result.WriteString(fmt.Sprintf("\t%-8s: %s\n", id, spec))
	}
	return result.String()
}

// from interface DataElemSpec
func (s *CompositeDataElemSpec) Id() string {
	return s.id
}

// from interface DataElemSpec
func (s *CompositeDataElemSpec) Name() string {
	return s.name
}

func (s *CompositeDataElemSpec) String() string {
	descrStr := util.Ellipsis(s.Descr, maxDescrDisplayLen)
	specsStrs := []string{}
	for _, spec := range s.ComponentSpecs {
		specsStrs = append(specsStrs, fmt.Sprintf("\t%s", spec.String()))
	}
	componentSpecsStr := strings.Join(specsStrs, "\n")
	return fmt.Sprintf(
		"Composite %s %s '%s'\n%s",
		s.id, s.name, descrStr, componentSpecsStr)
}

type CompositeDataElemSpecMap map[string]*CompositeDataElemSpec
type CompositeDataElemSpecs []*CompositeDataElemSpec

// from sort.Interface
func (m CompositeDataElemSpecs) Len() int {
	return len(m)
}

// from sort.Interface
func (m CompositeDataElemSpecs) Less(i, j int) bool {
	return m[i].id < m[j].id
}

// from sort.Interface
func (m CompositeDataElemSpecs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

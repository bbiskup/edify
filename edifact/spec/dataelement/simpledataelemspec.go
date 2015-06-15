package dataelement

import (
	"bytes"
	"fmt"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
	"github.com/bbiskup/edify/edifact/util"
)

// DataElem specification
type SimpleDataElemSpec struct {
	id         string
	name       string
	Descr      string
	Repr       *Repr
	CodesSpecs *csp.CodesSpec
}

func NewSimpleDataElemSpec(id string, name string, descr string, repr *Repr, codes *csp.CodesSpec) (*SimpleDataElemSpec, error) {
	err := util.CheckNotNil(id, name, descr, repr, codes)
	if err != nil {
		return nil, err
	}
	return &SimpleDataElemSpec{
		id:         id,
		name:       name,
		Descr:      descr,
		Repr:       repr,
		CodesSpecs: codes,
	}, nil
}

func (s *SimpleDataElemSpec) String() string {
	return fmt.Sprintf("SimpleDataElemSpec: %s '%s' [%s]", s.id, s.name, s.Repr)
}

// from interface DataElemSpec
func (s *SimpleDataElemSpec) Id() string {
	return s.id
}

// from interface DataElemSpec
func (s *SimpleDataElemSpec) Name() string {
	return s.name
}

func (sm SimpleDataElemSpecMap) String() string {
	var result bytes.Buffer
	result.WriteString("SimpleDataElemSpecMap\n")
	for id, spec := range sm {
		result.WriteString(fmt.Sprintf("\t%s: %s\n", id, spec))
	}
	return result.String()
}

type SimpleDataElemSpecMap map[string]*SimpleDataElemSpec
type SimpleDataElemSpecs []*SimpleDataElemSpec

// from sort.Interface
func (m SimpleDataElemSpecs) Len() int {
	return len(m)
}

// from sort.Interface
func (m SimpleDataElemSpecs) Less(i, j int) bool {
	return m[i].id < m[j].id
}

// from sort.Interface
func (m SimpleDataElemSpecs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

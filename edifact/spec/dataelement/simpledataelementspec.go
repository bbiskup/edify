package dataelement

import (
	"bytes"
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/codes"
	"github.com/bbiskup/edify/edifact/util"
)

// DataElement specification
type SimpleDataElementSpec struct {
	id         string
	name       string
	Descr      string
	Repr       *Repr
	CodesSpecs *codes.CodesSpec
}

func (s *SimpleDataElementSpec) String() string {
	return fmt.Sprintf("SimpleDataElementSpec: %s '%s' [%s]", s.id, s.name, s.Repr)
}

// from interface DataElementSpec
func (s *SimpleDataElementSpec) Id() string {
	return s.id
}

// from interface DataElementSpec
func (s *SimpleDataElementSpec) Name() string {
	return s.name
}

func NewSimpleDataElementSpec(id string, name string, descr string, repr *Repr, codes *codes.CodesSpec) (*SimpleDataElementSpec, error) {
	err := util.CheckNotNil(id, name, descr, repr, codes)
	if err != nil {
		return nil, err
	}
	return &SimpleDataElementSpec{
		id:         id,
		name:       name,
		Descr:      descr,
		Repr:       repr,
		CodesSpecs: codes,
	}, nil
}

type SimpleDataElementSpecMap map[string]*SimpleDataElementSpec

func (sm SimpleDataElementSpecMap) String() string {
	var result bytes.Buffer
	result.WriteString("SimpleDataElementSpecMap\n")
	for id, spec := range sm {
		result.WriteString(fmt.Sprintf("\t%s: %s\n", id, spec))
	}
	return result.String()
}

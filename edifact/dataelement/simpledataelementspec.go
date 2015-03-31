package dataelement

import (
	"fmt"
)

// DataElement specification
type SimpleDataElementSpec struct {
	Num   int32
	Name  string
	Descr string
	Repr  *Repr
}

func (s *SimpleDataElementSpec) String() string {
	return fmt.Sprintf("SimpleDataElementSpec: %d '%s' [%s]", s.Num, s.Name, s.Repr)
}

func NewSimpleDataElementSpec(num int32, name string, descr string, repr *Repr) *SimpleDataElementSpec {
	return &SimpleDataElementSpec{
		Num:   num,
		Name:  name,
		Descr: descr,
		Repr:  repr,
	}
}

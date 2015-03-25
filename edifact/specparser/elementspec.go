package specparser

import (
	"fmt"
)

// DataElement specification
type DataElementSpec struct {
	Num   int32
	Name  string
	Descr string
	Repr  string
}

func (s *DataElementSpec) String() string {
	return fmt.Sprintf("DataElementSpec: %d '%s' [%s]", s.Num, s.Name, s.Repr)
}

func NewDataElementSpec(num int32, name string, descr string, repr string) *DataElementSpec {
	return &DataElementSpec{
		Num:   num,
		Name:  name,
		Descr: descr,
		Repr:  repr,
	}
}

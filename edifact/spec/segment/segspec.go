package segment

import (
	"fmt"
)

// Seg specification
type SegSpec struct {
	Id               string
	Name             string
	Function         string
	SegDataElemSpecs []*SegDataElemSpec
}

func NewSegSpec(
	id string, name string, function string,
	segDataElemSpecs []*SegDataElemSpec) *SegSpec {

	return &SegSpec{
		Id:               id,
		Name:             name,
		Function:         function,
		SegDataElemSpecs: segDataElemSpecs,
	}
}

type SegSpecs []*SegSpec

func (s *SegSpec) String() string {
	return fmt.Sprintf(
		"Seg %s/%s (%d data elems)",
		s.Id, s.Name, len(s.SegDataElemSpecs))
}

// Determines how many component data elements must be present
// (trailing optional data elements may be omitted)
func (s *SegSpec) NumLeadingMandDataElems() int {
	numMandatory := 0
	for _, elemSpec := range s.SegDataElemSpecs {
		if !elemSpec.IsMandatory {
			break
		}
		numMandatory++
	}
	return numMandatory
}

// from sort.Interface
func (m SegSpecs) Len() int {
	return len(m)
}

// from sort.Interface
func (m SegSpecs) Less(i, j int) bool {
	return m[i].Id < m[j].Id
}

// from sort.Interface
func (m SegSpecs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

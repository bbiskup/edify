package segment

import (
	"fmt"
)

// Seg specification
type SegSpec struct {
	Id                      string
	Name                    string
	Function                string
	SegDataElemSpecs []*SegDataElemSpec
}

func (s *SegSpec) String() string {
	return fmt.Sprintf(
		"Seg %s/%s (%d data elems)",
		s.Id, s.Name, len(s.SegDataElemSpecs))
}

func NewSegSpec(
	id string, name string, function string,
	segmentDataElemSpecs []*SegDataElemSpec) *SegSpec {

	return &SegSpec{
		Id:                      id,
		Name:                    name,
		Function:                function,
		SegDataElemSpecs: segmentDataElemSpecs,
	}
}

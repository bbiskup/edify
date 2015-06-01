package segment

import (
	"fmt"
)

// Segment specification
type SegSpec struct {
	Id                      string
	Name                    string
	Function                string
	SegmentDataElemSpecs []*SegmentDataElemSpec
}

func (s *SegSpec) String() string {
	return fmt.Sprintf(
		"Segment %s/%s (%d data elems)",
		s.Id, s.Name, len(s.SegmentDataElemSpecs))
}

func NewSegSpec(
	id string, name string, function string,
	segmentDataElemSpecs []*SegmentDataElemSpec) *SegSpec {

	return &SegSpec{
		Id:                      id,
		Name:                    name,
		Function:                function,
		SegmentDataElemSpecs: segmentDataElemSpecs,
	}
}

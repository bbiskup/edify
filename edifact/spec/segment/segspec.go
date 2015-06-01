package segment

import (
	"fmt"
)

// Segment specification
type SegSpec struct {
	Id                      string
	Name                    string
	Function                string
	SegmentDataElementSpecs []*SegmentDataElementSpec
}

func (s *SegSpec) String() string {
	return fmt.Sprintf(
		"Segment %s/%s (%d data elems)",
		s.Id, s.Name, len(s.SegmentDataElementSpecs))
}

func NewSegSpec(
	id string, name string, function string,
	segmentDataElementSpecs []*SegmentDataElementSpec) *SegSpec {

	return &SegSpec{
		Id:                      id,
		Name:                    name,
		Function:                function,
		SegmentDataElementSpecs: segmentDataElementSpecs,
	}
}

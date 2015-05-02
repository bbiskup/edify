package segment

import (
	"fmt"
)

// Segment specification
type SegmentSpec struct {
	Id                      string
	Name                    string
	Function                string
	SegmentDataElementSpecs []*SegmentDataElementSpec
}

func (s *SegmentSpec) String() string {
	return fmt.Sprintf(
		"Segment %s/%s (%d data elems)",
		s.Id, s.Name, len(s.SegmentDataElementSpecs))
}

func NewSegmentSpec(
	id string, name string, function string,
	segmentDataElementSpecs []*SegmentDataElementSpec) *SegmentSpec {

	return &SegmentSpec{
		Id:                      id,
		Name:                    name,
		Function:                function,
		SegmentDataElementSpecs: segmentDataElementSpecs,
	}
}

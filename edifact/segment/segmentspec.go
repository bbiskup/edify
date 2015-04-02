package segment

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/dataelement"
)

// Segment specification
type SegmentSpec struct {
	Id               string
	Name             string
	Function         string
	DataElementSpecs []*dataelement.DataElementSpec
}

func (s *SegmentSpec) String() string {
	return fmt.Sprintf(
		"Segment %s/%s (%d data elems)",
		s.Id, s.Name, len(s.DataElementSpecs))
}

func NewSegmentSpec(id string, name string, function string, dataElementSpecs []*dataelement.DataElementSpec) *SegmentSpec {
	return &SegmentSpec{
		Id:               id,
		Name:             name,
		Function:         function,
		DataElementSpecs: dataElementSpecs,
	}
}

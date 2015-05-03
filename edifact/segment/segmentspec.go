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

type SegmentSpecMap map[string]*SegmentSpec

// Provides segment spec by Id
type SegmentSpecProvider interface {
	Get(id string) *SegmentSpec
	Len() int
}

// Regular implementation of SegmentSpecProvider for production
type SegmentSpecProviderImpl struct {
	segmentSpecs SegmentSpecMap
}

func (p *SegmentSpecProviderImpl) Get(id string) *SegmentSpec {
	return p.segmentSpecs[id]
}

func (p *SegmentSpecProviderImpl) Len() int {
	return len(p.segmentSpecs)
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

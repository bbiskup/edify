package segment

import (
	"fmt"
	"log"
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
	result := p.segmentSpecs[id]
	if result == nil {
		// e.g. UNH, UNT are not defined in UNCE specs, because they
		// are not part of the release cycle. Instead, they are defined
		// in part 1 of ISO9735 (file testdata/r1241.txt)
		log.Printf("######################## Missing segment spec: '%s'", id)
		return NewSegmentSpec(
			id, fmt.Sprintf("missing-%s", id),
			"dummy_function", nil)
	} else {
		return result
	}
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

package segment

import (
	"fmt"
	"log"
)

// Provides segment spec by Id
type SegmentSpecProvider interface {
	Get(id string) *SegmentSpec
	Len() int
}

type SegmentSpecMap map[string]*SegmentSpec

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

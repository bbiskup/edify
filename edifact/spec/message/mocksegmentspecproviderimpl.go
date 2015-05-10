package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
)

type MockSegmentSpecProviderImpl struct {
}

func (p *MockSegmentSpecProviderImpl) Get(id string) *segment.SegmentSpec {
	return segment.NewSegmentSpec(
		id, fmt.Sprintf("dummy_segment_spec-%s", id), "dummy_function", nil)
}

func (p *MockSegmentSpecProviderImpl) Len() int {
	// Dummy value; unused
	return 100
}

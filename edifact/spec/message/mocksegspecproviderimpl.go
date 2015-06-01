package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
)

type MockSegSpecProviderImpl struct {
}

func (p *MockSegSpecProviderImpl) Get(id string) *segment.SegSpec {
	return segment.NewSegSpec(
		id, fmt.Sprintf("dummy_segment_spec-%s", id), "dummy_function", nil)
}

func (p *MockSegSpecProviderImpl) Len() int {
	// Dummy value; unused
	return 100
}

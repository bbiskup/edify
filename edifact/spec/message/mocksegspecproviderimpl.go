package message

import (
	"fmt"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
)

type MockSegSpecProviderImpl struct {
}

func (p *MockSegSpecProviderImpl) Get(id string) *ssp.SegSpec {
	return ssp.NewSegSpec(
		id, fmt.Sprintf("dummy_segment_spec-%s", id), "dummy_function", nil)
}

func (p *MockSegSpecProviderImpl) Len() int {
	// Dummy value; unused
	return 100
}

func (p *MockSegSpecProviderImpl) Ids() []string {
	// Not useful; mock implementation will return a segment spec
	// for any ID
	panic("Not implemented")
}

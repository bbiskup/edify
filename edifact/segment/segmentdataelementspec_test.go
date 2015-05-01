package segment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// implementes dataelement.DataElementSpec
type DummyElem struct {
}

func (d *DummyElem) Id() string {
	return "dummy_id"
}

func (d *DummyElem) Name() string {
	return "dummy_name"
}

func TestSegmentDataElementSpec(t *testing.T) {
	e := NewSegmentDataElementSpec(&DummyElem{}, 3, true)
	const expected = "SegmentDataElem dummy_id 3x mand."
	assert.Equal(t, expected, e.String(), "String() incorrect")
}

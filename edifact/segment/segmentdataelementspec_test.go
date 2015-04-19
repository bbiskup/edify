package segment

import (
	"fmt"
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
	eStr := e.String()
	if eStr != expected {
		t.Errorf(fmt.Sprintf("Expected: %s, got: %s", expected, eStr))
	}
}

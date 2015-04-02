package segment

import (
	"testing"
)

func TestSegmentSpec(t *testing.T) {
	spec := NewSegmentSpec("ADR", "ADDRESS", "To specify an address.", nil)

	const expectedStr = "Segment ADR/ADDRESS (0 data elems)"
	specStr := spec.String()

	if specStr != expectedStr {
		t.Errorf("Expected: %s got: %s", expectedStr, specStr)
	}
}

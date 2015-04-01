package codes

import (
	"testing"
)

func TestCodeSpec(t *testing.T) {
	spec := NewCodeSpec("10", "testname", "testdescr")
	expected := "10 testname testdescr"
	if spec.String() != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, spec)
	}
}

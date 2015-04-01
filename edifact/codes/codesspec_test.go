package codes

import (
	"testing"
)

func TestCodesSpec(t *testing.T) {
	spec := NewCodesSpec(10, "codesspec_name", "codespec_descr", []*CodeSpec{
		NewCodeSpec("10", "codespec_name", "codespec_descr"),
	})
	expected := "10 codesspec_name codespe...\n\t10 codespec_name codespe..."
	if spec.String() != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, spec)
	}
}

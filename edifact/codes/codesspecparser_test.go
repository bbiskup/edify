package codes

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	parser := NewCodesSpecParser()
	const fileName = "../../testdata/UNCL.14B"
	specsMap, err := parser.ParseSpecFile(fileName)
	if err != nil {
		t.Fatalf("Failed to parse spec file %s: %s", fileName, err)
	}

	expectedLen := 272
	lenSpecsMap := len(specsMap)
	if lenSpecsMap != expectedLen {
		t.Errorf("Length mismatch; expected: %d; got: %d", expectedLen, lenSpecsMap)
	}
}

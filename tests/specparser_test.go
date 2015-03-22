package tests

import (
	sp "edifact_experiments/edifact/specparser"
	"testing"
)

func TestSpecParser(t *testing.T) {
	p := sp.NewDataElementSpecParser()
	_, err := p.ParseSpecFile("../testdata/EDED.14B_short")
	if err != nil {
		t.Fatalf("Parse error: %s", err)
	}
}

package tests

import (
	edi "edifact_experiments/edifact"
	"fmt"
	"testing"
)

var splitTests = []struct {
	input      string
	sep        rune
	escapeChar rune
	expected   []string
}{
	{
		input:      "abc",
		sep:        '+',
		escapeChar: '?',
		expected:   []string{"abc"},
	},

	{
		input:      "abc+def",
		sep:        '+',
		escapeChar: '?',
		expected:   []string{"abc", "def"},
	},

	{
		input:      "abc?+def",
		sep:        '+',
		escapeChar: '?',
		expected:   []string{"abc?+def"},
	},
}

func TestSplitEDIFACT(t *testing.T) {
	for _, s := range splitTests {
		res := edi.SplitEDIFACT(s.input, s.sep, s.escapeChar)

		if fmt.Sprintf("%#v", res) != fmt.Sprintf("%#v", s.expected) {
			t.Fatalf("Expected: %#v; got: %#c", s.expected, res)
		}
	}
}

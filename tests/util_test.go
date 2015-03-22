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
		input:      "",
		sep:        '+',
		escapeChar: '?',
		expected:   []string{},
	},
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

	{
		input:      "+abc+def+ghi++jkl+",
		sep:        '+',
		escapeChar: '?',
		expected:   []string{"", "abc", "def", "ghi", "", "jkl", ""},
	},
}

func TestSplitEDIFACT(t *testing.T) {
	for _, s := range splitTests {
		res := edi.SplitEDIFACT(s.input, s.sep, s.escapeChar)

		resStr := fmt.Sprintf("%#v", res)
		expectedStr := fmt.Sprintf("%#v", s.expected)
		if resStr != expectedStr {
			t.Fatalf("Expected: %#v; got: %#c", expectedStr, resStr)
		}
	}
}

var indentTests = []struct {
	str      string
	expected int
}{
	{"", 0},
	{" ", 1},
	{" a", 1},
	{"  ", 2},
	{"  a", 2},
	{"  a ", 2},
	{"  ab ", 2},
}

func TestGetIndent(t *testing.T) {
	for _, spec := range indentTests {
		res := edi.GetIndent(spec.str)
		if res != spec.expected {
			t.Fatalf("Failed for spec '%s': expected %d, got %d",
				spec.str, spec.expected, res)
		}
	}
}

package tests

import (
	edi "edifice/edifact"
	"fmt"
	"strings"
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

func BenchmarkSplitEDIFACT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		edi.SplitEDIFACT("abc+d?+ef+ghi", '+', '?')
	}
}

var getIndentTests = []struct {
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
	for _, spec := range getIndentTests {
		res := edi.GetIndent(spec.str)
		if res != spec.expected {
			t.Fatalf("Failed for spec '%s': expected %d, got %d",
				spec.str, spec.expected, res)
		}
	}
}

var splitByHangingIndentTests = []struct {
	lines    []string
	expected [][]string
}{
	{
		lines:    []string{},
		expected: [][]string{},
	},
	{
		lines: []string{
			"abc",
		},
		expected: [][]string{
			[]string{"abc"},
		},
	},
	{
		lines: []string{
			"abc",
			"def",
		},
		expected: [][]string{
			[]string{"abc"},
			[]string{"def"},
		},
	},
	{
		lines: []string{
			"abc",
			" def",
			" ghi",
		},
		expected: [][]string{
			[]string{"abc", " def", " ghi"},
		},
	},
	{
		lines: []string{
			"",
			"abc",
			" def",
			" ghi",
			"",
			"jkl",
			" mno",
			"",
		},
		expected: [][]string{
			[]string{"abc", " def", " ghi"},
			[]string{"jkl", " mno"},
		},
	},
	{
		lines: []string{
			"abc",
			" def",
			" ghi",
			"jkl",
			" mno",
			" pqr",
		},
		expected: [][]string{
			[]string{"abc", " def", " ghi"},
			[]string{"jkl", " mno", " pqr"},
		},
	},
}

func TestSplitByHangingIndent(t *testing.T) {
	for _, spec := range splitByHangingIndentTests {
		res := edi.SplitByHangingIndent(spec.lines, 0)

		expectedStr := fmt.Sprintf("%s", spec.expected)
		resStr := fmt.Sprintf("%s", res)

		if resStr != expectedStr {
			t.Fatalf("Failed for spec '%s': expected %s, got %s",
				spec.lines, expectedStr, res)
		}
	}
}

func BenchmarkSplitByHangingIndent(b *testing.B) {
	testStr := `

     8023  Freight and other charges description identifier        [B]

     Desc: Code identifying freight and other charges.

     Repr: an..17

     Note: 
           1 Use UN/ECE Recommendation No. 23: Freight costs and
           charges. If not applicable, use appropriate code in
           combination with 1131/3055.


	`
	testLines := strings.Split(testStr, "\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		edi.SplitByHangingIndent(testLines, 4)
	}
}

package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
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
		res := SplitEDIFACT(s.input, s.sep, s.escapeChar)
		assert.Equal(t, fmt.Sprintf("%#v", s.expected), fmt.Sprintf("%#v", res))
	}
}

func BenchmarkSplitEDIFACT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitEDIFACT("abc+d?+ef+ghi", '+', '?')
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
		assert.Equal(t, spec.expected, GetIndent(spec.str))
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
		res := SplitByHangingIndent(spec.lines, 0)

		assert.Equal(t, fmt.Sprintf("%s", spec.expected), fmt.Sprintf("%s", res))
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
		SplitByHangingIndent(testLines, 4)
	}
}

func TestCustBoolStr(t *testing.T) {
	assert.Equal(t, "yes", CustBoolStr(true, "yes", "no"))
	assert.Equal(t, "no", CustBoolStr(false, "yes", "no"))
}

var removeLeadingAndTrailingEmptyLinesSpecs = []struct {
	lines    []string
	expected []string
}{
	{
		lines:    []string{},
		expected: []string{},
	},
	{
		lines:    []string{"one"},
		expected: []string{"one"},
	},
	{
		lines:    []string{"one", "two"},
		expected: []string{"one", "two"},
	},
	{
		lines:    []string{"", "one", "two"},
		expected: []string{"one", "two"},
	},
	{
		lines:    []string{"", "one", "two", ""},
		expected: []string{"one", "two"},
	},
	{
		lines:    []string{"one", "two", ""},
		expected: []string{"one", "two"},
	},
}

func TestRemoveLeadingAndTrailingEmptyLines(t *testing.T) {
	for _, spec := range removeLeadingAndTrailingEmptyLinesSpecs {
		res := RemoveLeadingAndTrailingEmptyLines(spec.lines)
		assert.True(t, reflect.DeepEqual(res, spec.expected))
	}
}

var splitMultipleLinesByEmptyLinesSpecs = []struct {
	lines    []string
	expected [][]string
}{
	{
		lines:    []string{"one"},
		expected: [][]string{{"one"}},
	},
	{
		lines:    []string{"one", "two"},
		expected: [][]string{{"one", "two"}},
	},
	{
		lines:    []string{"one", "", "two"},
		expected: [][]string{{"one"}, {"two"}},
	},
	{
		lines:    []string{"one", "", "", "two"},
		expected: [][]string{{"one"}, {}, {"two"}},
	},
	{
		lines:    []string{"", "one", "two"},
		expected: [][]string{{}, {"one", "two"}},
	},
	{
		lines:    []string{"one", "two", ""},
		expected: [][]string{{"one", "two"}, {}},
	},
}

func TestSplitMultipleLinesByEmptyLinesSpecs(t *testing.T) {
	for _, spec := range splitMultipleLinesByEmptyLinesSpecs {
		res := SplitMultipleLinesByEmptyLines(spec.lines)
		assert.Equal(t, fmt.Sprintf("%s", spec.expected), fmt.Sprintf("%s", res))
	}
}

var ellipsisSpec = []struct {
	str      string
	maxLen   int
	expected string
}{
	{"", 0, ""},
	{"", 1, ""},
	{"one", 3, "one"},
	{"one", 4, "one"},
	{"one", 2, "..."},
	{"one", 0, "..."},
	{"onetwo", 5, "on..."},
	{"onetwo", 4, "o..."},
}

func TestEllipsis(t *testing.T) {
	for _, spec := range ellipsisSpec {
		res := Ellipsis(spec.str, spec.maxLen)
		assert.Equal(t, spec.expected, res)
	}
}

var joinByHangingIndentSpecs = []struct {
	lines          []string
	expected       []string
	baseIndent     int
	collapseSpaces bool
}{
	{
		[]string{},
		[]string{},
		0, true,
	},
	{
		[]string{"one"},
		[]string{"one"},
		0, true,
	},
	{
		[]string{"one", "two"},
		[]string{"one", "two"},
		0, true,
	},
	{
		[]string{"one", "  two"},
		[]string{"one  two"},
		0, false,
	},
	{
		[]string{"one", "  two"},
		[]string{"one two"},
		0, true,
	},
	{
		[]string{"one", "  two", "  three", "four"},
		[]string{"one two three", "four"},
		0, true,
	},
	{
		[]string{" one", "  two", "  three", " four"},
		[]string{"one two three", "four"},
		1, true,
	},
}

func TestJoinByHangingIndent(t *testing.T) {
	for _, spec := range joinByHangingIndentSpecs {
		res := JoinByHangingIndent(spec.lines, spec.baseIndent, spec.collapseSpaces)
		assert.Equal(t, fmt.Sprintf("%#v", spec.expected), fmt.Sprintf("%#v", res))
	}
}

var trimWhiteSpaceAndJoinSpecs = []struct {
	lines    []string
	joinStr  string
	expected string
}{
	{[]string{""}, " ", ""},
	{[]string{"", ""}, " ", " "},
	{[]string{"a"}, " ", "a"},
	{[]string{"a", "b"}, " ", "a b"},
	{[]string{"a", "b"}, "x", "axb"},
	{[]string{" a", "b"}, " ", "a b"},
	{[]string{" a ", " b "}, " ", "a b"},
	{[]string{"  a  ", "  b  "}, " ", "a b"},
	{[]string{"\ta\t", "\tb\t"}, " ", "a b"},
}

func TestTrimWhiteSpaceAndJoin(t *testing.T) {
	for _, spec := range trimWhiteSpaceAndJoinSpecs {
		res := TrimWhiteSpaceAndJoin(spec.lines, spec.joinStr)
		assert.Equal(t, spec.expected, res)
	}
}

var checkNotNilSpecs = []struct {
	values    []interface{}
	expectErr bool
}{
	{
		[]interface{}{},
		false,
	},
	{
		[]interface{}{1},
		false,
	},
	{
		[]interface{}{"x"},
		false,
	},
	{
		[]interface{}{nil},
		true,
	},
	{
		[]interface{}{1, nil},
		true,
	},
}

func TestCheckNotNil(t *testing.T) {
	for _, spec := range checkNotNilSpecs {
		err := CheckNotNil(spec.values...)
		if spec.expectErr {
			assert.NotNil(t, err)
		}
		if !spec.expectErr {
			assert.Nil(t, err)
		}
	}
}

func TestUnused(t *testing.T) {
	a := 1
	b := 2
	Unused(a, b)
}

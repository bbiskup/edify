package specparser

import (
	sp "edifice/edifact/specparser"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var reprSpec = []struct {
	reprStr  string
	expected *sp.Repr
}{
	{"a..1", sp.NewRepr(sp.Alpha, true, 1)},
	{"a..2", sp.NewRepr(sp.Alpha, true, 2)},
	{"a2", sp.NewRepr(sp.Alpha, false, 2)},
	{"an2", sp.NewRepr(sp.AlphaNum, false, 2)},
	{"an..2", sp.NewRepr(sp.AlphaNum, true, 2)},
	{"n3", sp.NewRepr(sp.Num, false, 3)},
	{"n..3", sp.NewRepr(sp.Num, true, 3)},
}

func TestParseRepr(t *testing.T) {
	for _, spec := range reprSpec {
		res, err := sp.ParseRepr(spec.reprStr)
		if err != nil {
			t.Fatalf("Parse error: %s", err)
		}
		if !reflect.DeepEqual(res, spec.expected) {
			t.Fatalf("Repr string: %s: expected: %#v, got: %#v",
				spec.reprStr, spec.expected, res)
		}
	}
}

var validationSpec = []struct {
	repr     *sp.Repr
	testStr  string
	expected bool
}{
	{sp.NewRepr(sp.Alpha, true, 1), "x", true},
	{sp.NewRepr(sp.Alpha, true, 1), "", true},
	{sp.NewRepr(sp.Alpha, true, 1), "xx", false},
	{sp.NewRepr(sp.Alpha, false, 2), "xx", true},
	{sp.NewRepr(sp.Alpha, false, 2), "x", false},

	{sp.NewRepr(sp.AlphaNum, true, 2), "x", true},
	{sp.NewRepr(sp.AlphaNum, true, 2), "xx", true},
	{sp.NewRepr(sp.AlphaNum, true, 3), "x2x", true},
	{sp.NewRepr(sp.AlphaNum, true, 3), "x2x2", false},
	{sp.NewRepr(sp.AlphaNum, true, 3), "", true},
	{sp.NewRepr(sp.AlphaNum, false, 4), "x2x2", true},

	{sp.NewRepr(sp.Num, true, 3), "123", true},
	{sp.NewRepr(sp.Num, true, 3), "123a", false},
	{sp.NewRepr(sp.Num, true, 3), "", true},
	{sp.NewRepr(sp.Num, false, 3), "123", true},
	{sp.NewRepr(sp.Num, false, 3), "12", false},
}

func TestValidateRepr(t *testing.T) {
	for _, spec := range validationSpec {
		res, err := spec.repr.Validate(spec.testStr)

		if spec.expected {
			if err != nil {
				t.Fatalf("Validation error: %s", err)
			}
		} else {
			if err == nil {
				t.Fatalf("Should get validation error")
			}
		}

		if !reflect.DeepEqual(res, spec.expected) {
			t.Fatalf("Test string: %s: expected: %#v, got: %#v",
				spec.testStr, spec.expected, res)
		}
	}
}

func BenchmarkValidateShortExpr(b *testing.B) {
	repr, err := sp.ParseRepr("an..10")
	if err != nil {
		b.Fatalf("Parse error: %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := repr.Validate("abcdef1234")
		if !res {
			b.Fatalf(fmt.Sprintf("Validation failed: %s", err))
		}
	}
}

func BenchmarkValidateLongExpr(b *testing.B) {
	repr, err := sp.ParseRepr("an..10000")
	if err != nil {
		b.Fatalf("Parse error: %s", err)
	}

	testStr := strings.Repeat("x", 10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := repr.Validate(testStr)
		if !res {
			b.Fatalf(fmt.Sprintf("Validation failed: %s", err))
		}
	}
}

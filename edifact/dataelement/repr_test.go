package dataelement

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var reprSpec = []struct {
	reprStr  string
	expected *Repr
}{
	{"a..1", NewRepr(Alpha, true, 1)},
	{"a..2", NewRepr(Alpha, true, 2)},
	{"a2", NewRepr(Alpha, false, 2)},
	{"an2", NewRepr(AlphaNum, false, 2)},
	{"an..2", NewRepr(AlphaNum, true, 2)},
	{"n3", NewRepr(Num, false, 3)},
	{"n..3", NewRepr(Num, true, 3)},
}

func TestParseRepr(t *testing.T) {
	for _, spec := range reprSpec {
		res, err := ParseRepr(spec.reprStr)
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
	repr     *Repr
	testStr  string
	expected bool
}{
	{NewRepr(Alpha, true, 1), "x", true},
	{NewRepr(Alpha, true, 1), "", true},
	{NewRepr(Alpha, true, 1), "xx", false},
	{NewRepr(Alpha, false, 2), "xx", true},
	{NewRepr(Alpha, false, 2), "x", false},

	{NewRepr(AlphaNum, true, 2), "x", true},
	{NewRepr(AlphaNum, true, 2), "xx", true},
	{NewRepr(AlphaNum, true, 3), "x2x", true},
	{NewRepr(AlphaNum, true, 3), "x2x2", false},
	{NewRepr(AlphaNum, true, 3), "", true},
	{NewRepr(AlphaNum, false, 4), "x2x2", true},

	{NewRepr(Num, true, 3), "123", true},
	{NewRepr(Num, true, 3), "123a", false},
	{NewRepr(Num, true, 3), "", true},
	{NewRepr(Num, false, 3), "123", true},
	{NewRepr(Num, false, 3), "12", false},
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
	repr, err := ParseRepr("an..10")
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
	repr, err := ParseRepr("an..10000")
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

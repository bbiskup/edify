package specparser

import (
	sp "edifice/edifact/specparser"
	"reflect"
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
	{sp.NewRepr(sp.AlphaNum, false, 4), "x2x2", true},

	{sp.NewRepr(sp.Num, true, 3), "123", true},
	{sp.NewRepr(sp.Num, true, 3), "123a", false},
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

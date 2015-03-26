package specparser

import (
	sp "edifice/edifact/specparser"
	"reflect"
	"testing"
)

var repSpec = []struct {
	reprStr  string
	expected *sp.Repr
}{
	{"a..1", sp.NewRepr(sp.Alpha, true, 1)},
	{"a..2", sp.NewRepr(sp.Alpha, true, 2)},
	{"a2", sp.NewRepr(sp.Alpha, false, 2)},
	{"an2", sp.NewRepr(sp.AlphaNum, false, 2)},
	{"an..2", sp.NewRepr(sp.AlphaNum, true, 2)},
}

func TestRep(t *testing.T) {
	for _, spec := range repSpec {
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

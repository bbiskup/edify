package tests

import (
	"fmt"
	de "github.com/bbiskup/edifice/edifact/dataelement"
	"testing"
)

func TestCompositeDataElementString(t *testing.T) {
	e1 := de.NewSimpleDataElementSpec(1, "elem1", "descr", de.NewRepr(de.Alpha, true, 3))
	elem := de.NewCompositeDataElementSpec("C817", "ADDRESS USAGE", []*de.SimpleDataElementSpec{
		e1,
	})
	expected := "C817 ADDRESS USAGE\n\tSimpleDataElementSpec: 1 'elem1' [a..3]"

	res := fmt.Sprintf("%s", elem)
	if res != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, res)
	}
}

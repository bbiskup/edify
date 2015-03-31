package tests

import (
	"fmt"
	de "github.com/bbiskup/edifice/edifact/dataelement"
	"testing"
)

func TestCompositeDataElementString(t *testing.T) {
	e1 := de.NewComponentDataElementSpec(1, true)
	elem := de.NewCompositeDataElementSpec("C817", "ADDRESS USAGE", 1, true, []*de.ComponentDataElementSpec{
		e1,
	})
	expected := "Composite C817 ADDRESS USAGE 1 (mandatory)\n\tComponent 1 (mandatory)"

	res := fmt.Sprintf("%s", elem)
	if res != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, res)
	}
}

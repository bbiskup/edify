package dataelement

import (
	"fmt"
	"testing"
)

func TestCompositeDataElementString(t *testing.T) {
	e1 := NewComponentDataElementSpec(1, true)
	elem := NewCompositeDataElementSpec("C817", "ADDRESS USAGE", "test description", []*ComponentDataElementSpec{
		e1,
	})
	expected := "Composite C817 ADDRESS USAGE 'test de...'\n\tComponent 1 (mand.)"

	res := fmt.Sprintf("%s", elem)
	if res != expected {
		t.Fatalf("Expected: '%s', got: '%s'", expected, res)
	}
}

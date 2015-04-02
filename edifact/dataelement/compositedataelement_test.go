package dataelement

import (
	"fmt"
	"testing"
)

func TestCompositeDataElementString(t *testing.T) {
	codesSpecMap := fixtureTextCodesSpecMap()
	simpleDataElemSpec, err := NewSimpleDataElementSpec(1, "name_1", "descr_1", NewRepr(AlphaNum, true, 3), codesSpecMap[1000])
	if err != nil {
		t.Errorf("Failed to create simple data element spec: %s", err)
	}

	e1 := NewComponentDataElementSpec(1, true, simpleDataElemSpec)
	elem := NewCompositeDataElementSpec("C817", "ADDRESS USAGE", "test description", []*ComponentDataElementSpec{
		e1,
	})
	expected := "Composite C817 ADDRESS USAGE 'test de...'\n\tComponent 1 (mand.)"

	res := fmt.Sprintf("%s", elem)
	if res != expected {
		t.Errorf("Expected: '%s', got: '%s'", expected, res)
	}
}

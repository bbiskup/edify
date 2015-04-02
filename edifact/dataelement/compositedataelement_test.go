package dataelement

import (
	"testing"
)

const expectedStr = `Composite C817 ADDRESS USAGE 'test de...'
	Component 1/name_1 @ 1 (mand.)`

func TestCompositeDataElementString(t *testing.T) {
	codesSpecMap := fixtureTextCodesSpecMap()
	simpleDataElemSpec, err := NewSimpleDataElementSpec("1", "name_1", "descr_1", NewRepr(AlphaNum, true, 3), codesSpecMap["1000"])
	if err != nil {
		t.Errorf("Failed to create simple data element spec: %s", err)
	}

	e1 := NewComponentDataElementSpec(1, true, simpleDataElemSpec)
	compositeDataElemSpec := NewCompositeDataElementSpec("C817", "ADDRESS USAGE", "test description", []*ComponentDataElementSpec{
		e1,
	})

	if compositeDataElemSpec.Id() != "C817" || compositeDataElemSpec.id != "C817" {
		t.Errorf("incorrect Id")
	}

	if compositeDataElemSpec.Name() != "ADDRESS USAGE" || compositeDataElemSpec.name != "ADDRESS USAGE" {
		t.Errorf("incorrect Name()")
	}

	specStr := compositeDataElemSpec.String()
	if compositeDataElemSpec.String() != expectedStr {
		t.Errorf("incorrect String(): %s", specStr)
	}
}

package dataelement

import (
	"testing"
)

func TestSimpleDataElementSpec(t *testing.T) {
	codesSpecMap := fixtureTextCodesSpecMap()
	spec, err := NewSimpleDataElementSpec("1", "name_1", "descr_1", NewRepr(AlphaNum, true, 3), codesSpecMap["1000"])
	if err != nil {
		t.Errorf("Failed to create simple data element spec: %s", err)
	}

	if spec.Id() != "1" || spec.id != "1" {
		t.Errorf("incorrect Id")
	}

	if spec.Name() != "name_1" || spec.name != "name_1" {
		t.Errorf("incorrect Name()")
	}

	const expectedStr = "SimpleDataElementSpec: 1 'name_1' [an..3]"
	specStr := spec.String()
	if spec.String() != expectedStr {
		t.Errorf("incorrect String(): %s", specStr)
	}
}

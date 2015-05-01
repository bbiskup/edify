package dataelement

import (
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "C817", compositeDataElemSpec.Id(), "Incorrect Id()")
	assert.Equal(t, "C817", compositeDataElemSpec.id, "Incorrect id")
	assert.Equal(t, "ADDRESS USAGE", compositeDataElemSpec.Name(), "Incorrect Name()")
	assert.Equal(t, "ADDRESS USAGE", compositeDataElemSpec.name, "Incorrect name")
	assert.Equal(t, expectedStr, compositeDataElemSpec.String(), "IncorrectString")
}

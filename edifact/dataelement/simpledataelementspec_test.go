package dataelement

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElementSpec(t *testing.T) {
	codesSpecMap := fixtureTextCodesSpecMap()
	spec, err := NewSimpleDataElementSpec("1", "name_1", "descr_1", NewRepr(AlphaNum, true, 3), codesSpecMap["1000"])
	assert.Nil(t, err, "Failed to create simple data element spec")
	assert.Equal(t, "1", spec.Id(), "Incorrect Id()")
	assert.Equal(t, "1", spec.id, "Incorrect id")
	assert.Equal(t, "name_1", spec.Name(), "Incorrect Name()")
	assert.Equal(t, "name_1", spec.name, "Incorrect name")

	const expectedStr = "SimpleDataElementSpec: 1 'name_1' [an..3]"
	assert.Equal(t, "SimpleDataElementSpec: 1 'name_1' [an..3]", spec.String(), "Incorrect String()")
}

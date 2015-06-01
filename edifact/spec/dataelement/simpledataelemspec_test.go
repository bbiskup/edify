package dataelement

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElemSpec(t *testing.T) {
	codesSpecMap := fixtureTextCodesSpecMap()
	spec, err := NewSimpleDataElemSpec("1", "name_1", "descr_1", NewRepr(AlphaNum, true, 3), codesSpecMap["1000"])
	assert.Nil(t, err, "Failed to create simple data element spec")
	assert.Equal(t, "1", spec.Id())
	assert.Equal(t, "1", spec.id)
	assert.Equal(t, "name_1", spec.Name())
	assert.Equal(t, "name_1", spec.name)

	const expectedStr = "SimpleDataElemSpec: 1 'name_1' [an..3]"
	assert.Equal(t, "SimpleDataElemSpec: 1 'name_1' [an..3]", spec.String())
}

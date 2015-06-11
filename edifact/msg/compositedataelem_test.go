package msg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompositeDataElem(t *testing.T) {
	elem := NewCompositeDataElem("ABC", NewSimpleDataElem("DEF", "GHI"))
	assert.Equal(t, "ABC", elem.Id())
	assert.Equal(t, 1, len(elem.SimpleDataElems))
	assert.Equal(t, "CompositeDataElem ABC (1 simple data elems)", elem.String())

	simpleDataElemById, err := elem.GetSimpleDataElemById("DEF")
	require.Nil(t, err)
	assert.Equal(t, "DEF", simpleDataElemById.Id())
	assert.Equal(t, "GHI", simpleDataElemById.Value)
}

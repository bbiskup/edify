package msg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSeg(t *testing.T) {
	seg := NewSeg("ABC")
	seg.AddDataElem(NewCompositeDataElem("DEF"))
	seg.AddDataElem(NewSimpleDataElem("GHI", "JKL"))
	assert.Equal(t, "ABC", seg.Id())

	compositeDataElemById, err := seg.GetCompositeDataElemById("DEF")
	require.Nil(t, err)
	assert.Equal(t, "DEF", compositeDataElemById.Id())

	simpleDataElemById, err := seg.GetSimpleDataElemById("GHI")
	require.Nil(t, err)
	assert.Equal(t, "GHI", simpleDataElemById.Id())
	assert.Equal(t, "JKL", simpleDataElemById.Value)

	assert.Equal(t, 2, len(seg.DataElems))
}

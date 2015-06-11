package msg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSeg(t *testing.T) {
	seg := NewSeg("ABC")
	seg.AddDataElem(NewCompositeDataElem("DEF"))
	assert.Equal(t, "ABC", seg.Id())

	dataElemById, err := seg.GetDataElemById("DEF")
	require.Nil(t, err)
	assert.Equal(t, "DEF", dataElemById.Id())
	assert.Equal(t, 1, len(seg.DataElems))
}

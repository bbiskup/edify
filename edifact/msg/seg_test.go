package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeg(t *testing.T) {
	seg := NewSeg("ABC")
	assert.Equal(t, "ABC", seg.Id())
	assert.Equal(t, 0, len(seg.DataElems))
}

package rawmsg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringEmpty(t *testing.T) {
	msg := NewRawMsg("ABC", []*RawSeg{})
	assert.Equal(t, "ABC", msg.String())
}

func TestStringWithSegs(t *testing.T) {
	segment1 := NewRawSeg("XYZ")
	msg := NewRawMsg("ABC", []*RawSeg{segment1})
	assert.Equal(t, 1, len(msg.RawSegs))
	assert.Equal(t, "ABC\n\tXYZ\n", msg.String())
}

func TestSegIds(t *testing.T) {
	rawSeg1 := NewRawSeg("DEF")
	rawSeg2 := NewRawSeg("GHI")
	msg := NewRawMsg("ABC", []*RawSeg{rawSeg1, rawSeg2})
	rawSegIds := msg.RawSegIds()
	assert.Equal(t, []string{"DEF", "GHI"}, rawSegIds)
}

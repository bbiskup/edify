package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringEmpty(t *testing.T) {
	msg := NewRawMessage("ABC", []*Seg{})
	assert.Equal(t, "ABC", msg.String())
}

func TestStringWithSegs(t *testing.T) {
	segment1 := NewSeg("XYZ")
	msg := NewRawMessage("ABC", []*Seg{segment1})
	assert.Equal(t, 1, len(msg.Segs))
	assert.Equal(t, "ABC\n\tXYZ\n", msg.String())
}

func TestSegIds(t *testing.T) {
	segment1 := NewSeg("DEF")
	segment2 := NewSeg("GHI")
	msg := NewRawMessage("ABC", []*Seg{segment1, segment2})
	segmentIds := msg.SegIds()
	assert.Equal(t, []string{"DEF", "GHI"}, segmentIds)
}

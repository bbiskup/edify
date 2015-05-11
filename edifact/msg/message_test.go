package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringEmpty(t *testing.T) {
	msg := NewMessage("ABC", []*Segment{})
	assert.Equal(t, "ABC\n", msg.String())
}

func TestStringWithSegments(t *testing.T) {
	segment1 := NewSegment("XYZ")
	msg := NewMessage("ABC", []*Segment{segment1})
	assert.Equal(t, 1, len(msg.Segments))
	assert.Equal(t, "ABC\n\tXYZ\n\n", msg.String())
}

func TestSegmentIds(t *testing.T) {
	segment1 := NewSegment("DEF")
	segment2 := NewSegment("GHI")
	msg := NewMessage("ABC", []*Segment{segment1, segment2})
	segmentIds := msg.SegmentIds()
	assert.Equal(t, []string{"DEF", "GHI"}, segmentIds)
}
package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringEmpty(t *testing.T) {
	msg := NewRawMessage("ABC", []*Segment{})
	assert.Equal(t, "ABC", msg.String())
}

func TestStringWithSegments(t *testing.T) {
	segment1 := NewSegment("XYZ")
	msg := NewRawMessage("ABC", []*Segment{segment1})
	assert.Equal(t, 1, len(msg.Segments))
	assert.Equal(t, "ABC\n\tXYZ\n\n", msg.String())
}

func TestSegmentIds(t *testing.T) {
	segment1 := NewSegment("DEF")
	segment2 := NewSegment("GHI")
	msg := NewRawMessage("ABC", []*Segment{segment1, segment2})
	segmentIds := msg.SegmentIds()
	assert.Equal(t, []string{"DEF", "GHI"}, segmentIds)
}

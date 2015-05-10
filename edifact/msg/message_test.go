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

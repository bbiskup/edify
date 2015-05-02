package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageSpecGroupPart(t *testing.T) {
	part := NewMessageSpecSegmentGroupPart(
		"testgroup", []MessageSpecPart{}, 5, true)
	assert.Equal(t, "Segment group testgroup 5 mand. (0 children)", part.String())
	assert.Equal(t, part.Name, "testgroup")
	assert.Equal(t, 5, part.MaxCount())
	assert.Equal(t, true, part.IsGroup())
	assert.Equal(t, 0, part.Count())
}

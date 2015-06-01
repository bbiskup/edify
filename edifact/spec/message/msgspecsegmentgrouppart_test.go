package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsgSpecGroupPart(t *testing.T) {
	part := NewMsgSpecSegmentGroupPart(
		"testgroup", []MsgSpecPart{}, 5, true, nil)
	assert.Equal(t, "Segment group testgroup 5 mand. (0 children)", part.String())
	assert.Equal(t, part.Name(), "testgroup")
	assert.Equal(t, 5, part.MaxCount())
	assert.Equal(t, true, part.IsGroup())
	assert.Equal(t, 0, part.Count())
	assert.Equal(t, nil, part.Parent())
}

func TestAppend(t *testing.T) {
	part1 := NewMsgSpecSegmentGroupPart(
		"testgroup1", []MsgSpecPart{}, 5, true, nil)
	assert.Equal(t, 0, part1.Count())
	part2 := NewMsgSpecSegmentGroupPart(
		"testgroup2", []MsgSpecPart{}, 5, true, nil)

	part1.Append(part2)
	assert.Equal(t, 1, part1.Count())
}

func TestParent(t *testing.T) {
	part1 := NewMsgSpecSegmentGroupPart(
		"testgroup1", []MsgSpecPart{}, 5, true, nil)
	part2 := NewMsgSpecSegmentGroupPart(
		"testgroup2", []MsgSpecPart{}, 5, true, part1)
	assert.Equal(t, part1, part2.Parent())
}

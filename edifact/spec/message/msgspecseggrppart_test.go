package message

import (
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsgSpecGroupPart(t *testing.T) {
	part := NewMsgSpecSegGrpPart(
		"testgroup", []MsgSpecPart{}, 5, true, nil)
	assert.Equal(t, "Segment group testgroup 5 mand. (0 children)", part.String())
	assert.Equal(t, part.Name(), "testgroup")
	assert.Equal(t, 5, part.MaxCount())
	assert.Equal(t, true, part.IsGroup())
	assert.Equal(t, 0, part.Count())
	assert.Equal(t, nil, part.Parent())
}

func TestAppend(t *testing.T) {
	part1 := NewMsgSpecSegGrpPart(
		"testgroup1", []MsgSpecPart{}, 5, true, nil)
	assert.Equal(t, 0, part1.Count())
	part2 := NewMsgSpecSegGrpPart(
		"testgroup2", []MsgSpecPart{}, 5, true, nil)

	part1.Append(part2)
	assert.Equal(t, 1, part1.Count())
}

func TestParent(t *testing.T) {
	part1 := NewMsgSpecSegGrpPart(
		"testgroup1", []MsgSpecPart{}, 5, true, nil)
	part2 := NewMsgSpecSegGrpPart(
		"testgroup2", []MsgSpecPart{}, 5, true, part1)
	assert.Equal(t, part1, part2.Parent())
}

func TestFindSegGrpSpecSimple(t *testing.T) {
	part := NewMsgSpecSegGrpPart(
		"testgroup",
		[]MsgSpecPart{
			NewMsgSpecSegPart(ssp.NewSegSpec("ABC", "ABC_NAME", "ABC_fun", nil), 1, true, nil),
			NewMsgSpecSegGrpPart("Group_1", nil, 4, true, nil),
		},
		5, true, nil)

	grp, err := part.FindSegGrpSpec("Group_1")
	assert.Nil(t, err)
	assert.Equal(t, "Group_1", grp.Name())
	assert.Equal(t, 4, grp.MaxCount())
}

func TestFindSegGrpSpecNested(t *testing.T) {
	part := NewMsgSpecSegGrpPart(
		"testgroup",
		[]MsgSpecPart{
			NewMsgSpecSegPart(ssp.NewSegSpec("ABC", "ABC_NAME", "ABC_fun", nil), 1, true, nil),
			NewMsgSpecSegGrpPart("Group_1",
				[]MsgSpecPart{
					NewMsgSpecSegGrpPart("Group_2", nil, 3, true, nil),
				},
				1, true, nil),
		},
		5, true, nil)

	grp, err := part.FindSegGrpSpec("Group_2")
	assert.Nil(t, err)
	assert.Equal(t, "Group_2", grp.Name())
	assert.Equal(t, 3, grp.MaxCount())
}

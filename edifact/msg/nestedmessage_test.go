package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEmptyNestedMsg() *NestedMessage {
	return NewNestedMessage("testname", []SegmentOrGroup{})
}

func getNestedMsgWithParts() *NestedMessage {
	return NewNestedMessage("testname", []SegmentOrGroup{
		NewSegment("ABC"),
		NewSegment("DEF"),
	})
}

func getNestedMsgWithGroupPart() *NestedMessage {
	return NewNestedMessage("testname", []SegmentOrGroup{
		NewSegment("ABC"),
		NewSegmentGroup("group_1", []SegmentOrGroup{
			NewSegment("DEF"),
			NewSegment("GHI"),
			NewSegmentGroup("group_2", []SegmentOrGroup{
				NewSegment("JKL"),
			}),
		}),
		NewSegment("MNO"),
	})
}

func TestStringEmptyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	assert.Equal(t, "NestedMessage testname (0 1st-level parts)", msg.String())
}

func SegGroupDumpEmptyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	assert.Equal(t, "<no msg parts>", msg.SegGroupDump())
}

func TestStringMsgWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, "NestedMessage testname (2 1st-level parts)", msg.String())
}

func TestSegGroupDumpWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, "ABC\nDEF\n", msg.SegGroupDump())
}

func TestSegGroupDumpWithGroupParts(t *testing.T) {
	msg := getNestedMsgWithGroupPart()
	assert.Equal(
		t,
		"ABC\ngroup_1\n\tDEF\n\tGHI\n\tgroup_2\n\t\tJKL\nMNO\n",
		msg.SegGroupDump())
}

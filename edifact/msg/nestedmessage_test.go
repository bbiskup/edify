package msg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEmptyNestedMsg() *NestedMessage {
	return NewNestedMessage("testname", []RepeatMsgPart{})
}

func getNestedMsgWithParts() *NestedMessage {
	return NewNestedMessage("testname", []RepeatMsgPart{
		NewRepeatSegment(NewSegment("ABC")),
		NewRepeatSegment(NewSegment("DEF")),
	})
}

func getNestedMsgWithGroupPart() *NestedMessage {
	return NewNestedMessage(
		"testname",
		[]RepeatMsgPart{
			NewRepeatSegment(
				NewSegment("ABC"),
				NewSegment("ABC"),
			),
			NewRepeatSegmentGroup(
				NewSegmentGroup("group_1", []RepeatMsgPart{
					NewRepeatSegment(NewSegment("DEF")),
					NewRepeatSegment(NewSegment("GHI")),
					NewRepeatSegmentGroup(NewSegmentGroup("group_2",
						[]RepeatMsgPart{
							NewRepeatSegment(
								NewSegment("JKL"),
							)})),

					NewRepeatSegment(NewSegment("MNO")),
				}),
			),
		})
}

func TestStringEmptyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	assert.Equal(t, "NestedMessage testname (0 1st-level parts)", msg.String())
}

func SegGroupyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	dump := msg.Dump(0)
	t.Logf("Dump:\n%s\n", dump)
	assert.Equal(t, "<no msg parts>", dump)
}

func TestStringMsgWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, "NestedMessage testname (2 1st-level parts)", msg.String())
}

func TestDumpWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	dump := msg.Dump(0)
	t.Logf("Dump:\n%s\n", dump)
	//assert.Equal(t, "ABC\nDEF\n", dump)
}

func TestDumpWithGroupParts(t *testing.T) {
	msg := getNestedMsgWithGroupPart()
	dump := msg.Dump(0)
	fmt.Printf("Dump:\n%s\n", dump)
	//assert.Equal(
	//	t,
	//	"ABC\ngroup_1\n\tDEF\n\tGHI\n\tgroup_2\n\t\tJKL\nMNO\n",
	//	dump)
}

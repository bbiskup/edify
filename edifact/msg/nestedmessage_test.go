package msg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEmptyNestedMsg() *NestedMessage {
	return NewNestedMessage("testname", []RepeatMsgPart{})
}

const expectedDumpNestedMsgWithParts = `Message testname
  [0] ABC
  [0] DEF
`

func getNestedMsgWithParts() *NestedMessage {
	return NewNestedMessage("testname", []RepeatMsgPart{
		NewRepeatSegment(NewSegment("ABC")),
		NewRepeatSegment(NewSegment("DEF")),
	})
}

const expectedDumpNestedMsgWithGroupPart = `Message testname
  [0] ABC
  [1] ABC
  [0] Group group_1
        [0] DEF
        [0] GHI
        [0] Group group_2
            [0] JKL
        [0] MNO
`

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
	dump := msg.Dump()
	t.Logf("Dump:\n%s\n", dump)
	assert.Equal(t, "<no msg parts>", dump)
}

func TestStringMsgWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, "NestedMessage testname (2 1st-level parts)", msg.String())
}

func TestDumpWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	dump := msg.Dump()
	t.Logf("Dump:\n%s\n", dump)
	assert.Equal(t, expectedDumpNestedMsgWithParts, dump)
}

func TestDumpWithGroupParts(t *testing.T) {
	msg := getNestedMsgWithGroupPart()
	dump := msg.Dump()
	fmt.Printf("Dump:\n%s\n", dump)
	assert.Equal(t, expectedDumpNestedMsgWithGroupPart, dump)
}

func TestAppend(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, 2, msg.Count())
	msg.AppendPart(NewRepeatSegment(NewSegment("ABC")))
	assert.Equal(t, 3, msg.Count())
	assert.Equal(t, "ABC", msg.GetPart(2).Id())
}

func BenchmarkDumpWithGroupParts(b *testing.B) {
	msg := getNestedMsgWithGroupPart()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dump := msg.Dump()
		assert.NotNil(b, dump)
	}
}

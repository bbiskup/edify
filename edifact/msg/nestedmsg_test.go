package msg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEmptyNestedMsg() *NestedMsg {
	return NewNestedMsg("testname")
}

const expectedDumpNestedMsgWithParts = `Message testname
  RepSegGrp
    [0] Group _toplevel
      RepSeg
        [0] ABC
      RepSeg
        [0] DEF

`

func getNestedMsgWithParts() *NestedMsg {
	return NewNestedMsg("testname",
		NewRepSeg(NewSeg("ABC")),
		NewRepSeg(NewSeg("DEF")),
	)
}

const expectedDumpNestedMsgWithGroupPart = `Message testname
  RepSegGrp
    [0] Group _toplevel
      RepSeg
        [0] ABC
        [1] ABC
      RepSegGrp
        [0] Group group_1
          RepSeg
            [0] DEF
          RepSeg
            [0] GHI
          RepSegGrp
            [0] Group group_2
              RepSeg
                [0] JKL
          RepSeg
            [0] MNO

`

func getNestedMsgWithGroupPart() *NestedMsg {
	return NewNestedMsg(
		"testname",
		NewRepSeg(
			NewSeg("ABC"),
			NewSeg("ABC"),
		),
		NewRepSegGrp(
			NewSegGrp("group_1",
				NewRepSeg(NewSeg("DEF")),
				NewRepSeg(NewSeg("GHI")),
				NewRepSegGrp(NewSegGrp("group_2",
					NewRepSeg(
						NewSeg("JKL"),
					))),
				NewRepSeg(NewSeg("MNO")),
			),
		),
	)
}

func TestStringEmptyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	assert.Equal(t, "NestedMsg testname (0 1st-level parts)", msg.String())
}

func SegGroupyMsg(t *testing.T) {
	msg := getEmptyNestedMsg()
	dump := msg.Dump()
	t.Logf("Dump:\n%s\n", dump)
	assert.Equal(t, "<no msg parts>", dump)
}

func TestStringMsgWithParts(t *testing.T) {
	msg := getNestedMsgWithParts()
	assert.Equal(t, "NestedMsg testname (2 1st-level parts)", msg.String())
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
	grp1 := NewSegGrp("Group_1", NewRepSeg(NewSeg("GHI")))
	msg.GetTopLevelGroup().AppendRepSegGrp(NewRepSegGrp(grp1))
	assert.Equal(t, 3, msg.Count())
	assert.Equal(t, "ABC", msg.GetTopLevelGroup().GetPart(0).Id())
	assert.Equal(t, 1, msg.GetTopLevelGroup().GetPart(0).Count())

	grp1_fetched := msg.GetTopLevelGroup().GetPart(msg.Count() - 1).(*RepSegGrp)
	assert.Equal(t, "Group_1", grp1_fetched.Id())

	ghi_fetched := grp1_fetched.GetSegGrp(0).GetPart(0)
	assert.Equal(t, "GHI", ghi_fetched.Id())
}

func TestSegGroupMap(t *testing.T) {
	msg := getNestedMsgWithGroupPart()
	grp := msg.TopLevelRepGrp.GetSegGrp(0)
	assert.True(t, grp.Contains("ABC"))
	grp.AppendRepSeg(NewRepSeg(NewSeg("XYZ")))
	assert.True(t, grp.Contains("XYZ"))

	grp1 := grp.GetPartByKey("group_1").(*RepSegGrp).GetSegGrp(0)
	assert.True(t, grp1.Contains("GHI"))
	assert.False(t, grp1.Contains("GHI2"))
}

func BenchmarkDumpWithGroupParts(b *testing.B) {
	msg := getNestedMsgWithGroupPart()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dump := msg.Dump()
		assert.NotNil(b, dump)
	}
}

func BenchmarkCreateNestedMsg(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := getNestedMsgWithGroupPart()
		assert.NotNil(b, msg)
	}
}

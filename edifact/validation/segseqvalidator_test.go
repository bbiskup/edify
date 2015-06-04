package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var authorSegSeqSpec = []struct {
	descr             string
	segmentIDs        []string
	expectError       bool
	errorKind         SegSeqErrKind
	validateNestedMsg func(t *testing.T, nestedMsg *msg.NestedMsg)
}{
	{
		"No segments at all",
		[]string{}, true, noSegs, nil,
	},

	{
		"Missing mandatory segments",
		[]string{
			"UNH", // no BGM
		}, true, missingMandatorySeg, nil,
	},

	{
		"Max. repeat count of mandatory segment exeeded",
		[]string{
			"UNH", "UNH", // max. repeat count is 1
		}, true, maxSegRepeatCountExceeded, nil,
	},

	{
		"Max. repeat count of optional segment exeeded",
		[]string{
			"UNH", "BGM", "DTM" /* max. repeat count is 1 */, "DTM",
		}, true, maxSegRepeatCountExceeded, nil,
	},

	{"Optional segment in incorrect position",
		[]string{
			"UNH",
			"DTM" /* Should appear after BGM */, "BGM", "UNT",
		}, true, missingMandatorySeg, nil,
	},

	{"Optional segment in incorrect position",
		[]string{
			"DTM", "UNH", "BGM", "UNT",
		}, true, missingMandatorySeg, nil,
	},

	{"Missing mandatory group 4",
		[]string{
			"UNH", "BGM", "DTM" /* optional */, "UNT",
		}, true, missingMandatorySeg, nil,
	},

	{"minimal message (only mandatory segments)",
		[]string{
			"UNH", "BGM" /* Group 4 */, "LIN",

			"UNT",
		}, false, "",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			assert.Equal(t, 4, nestedMsg.Count())
			assert.Equal(t, 4, nestedMsg.GetTopLevelGrp().Count())
			assert.Equal(t, "UNH", nestedMsg.GetTopLevelGrp().GetPart(0).Id())
			assert.Equal(t, "UNT", nestedMsg.GetTopLevelGrp().GetPart(3).Id())
		},
	},

	{
		"Mostly mandatory",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS", // both conditional
			// Group 4
			"LIN",
			"UNT",
		}, false, "",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			assert.Equal(t, 6, nestedMsg.Count())
			assert.Equal(t, 6, nestedMsg.GetTopLevelGrp().Count())
			assert.Equal(t, "UNH", nestedMsg.GetTopLevelGrp().GetPart(0).Id())
			assert.Equal(t, "DTM", nestedMsg.GetTopLevelGrp().GetPart(2).Id())
			assert.Equal(t, "UNT", nestedMsg.GetTopLevelGrp().GetPart(5).Id())
		},
	},

	{
		"Mostly mandatory; one conditional group",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",

			// Group 2
			"FII", "CTA", "COM",

			// Group 4
			"LIN",

			"UNT",
		}, false, "",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			assert.Equal(t, 7, nestedMsg.Count())
			assert.Equal(t, 7, nestedMsg.GetTopLevelGrp().Count())
			assert.Equal(t, "UNH", nestedMsg.GetTopLevelGrp().GetPart(0).Id())
			assert.Equal(t, "Group_2", nestedMsg.GetTopLevelGrp().GetPart(4).Id())
			assert.Equal(t, "Group_4", nestedMsg.GetTopLevelGrp().GetPart(5).Id())
		},
	},

	{
		"Nested conditional group",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",

			// Group 4
			"LIN",

			// Group 7 (nested in group 4)
			"FII", "CTA", "COM",

			"UNT",
		}, false, "",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			assert.Equal(t, 6, nestedMsg.Count())
			assert.Equal(t, 6, nestedMsg.GetTopLevelGrp().Count())

			unh := nestedMsg.GetTopLevelGrp().GetPart(0)
			assert.Equal(t, "UNH", unh.Id())
			assert.Equal(t, 1, unh.Count())
			group_4 := nestedMsg.GetTopLevelGrp().GetPart(4).(*msg.RepSegGrp)
			assert.Equal(t, "Group_4", group_4.Id())
			assert.Equal(t, 1, group_4.Count())
			assert.Equal(t, 2, group_4.GetSegGrp(0).Count())
			group_7 := group_4.GetSegGrp(0).GetPart(1).(*msg.RepSegGrp)
			assert.Equal(t, "Group_7", group_7.Id())
			assert.Equal(t, 1, group_7.Count())
			assert.Equal(t, 3, group_7.GetSegGrp(0).Count())
			assert.Equal(t, "UNT", nestedMsg.GetTopLevelGrp().GetPart(5).Id())
		},
	},

	{
		"Some repeat counts > 1",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 4
			"LIN", "LIN", "LIN", "LIN",
			// Group 7
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",

			"UNT",
		}, false, "",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			assert.Equal(t, 6, nestedMsg.Count())
			group_4 := nestedMsg.GetTopLevelGrp().GetPart(4).(*msg.RepSegGrp)
			assert.Equal(t, 4, group_4.Count())
			lin := group_4.GetSegGrp(0).GetPart(0).(*msg.RepSeg)
			assert.Equal(t, "LIN", lin.Id())

		},
	},

	{
		"group 7 repeated too often",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 4
			"LIN", "LIN", "LIN", "LIN",
			// Group 7
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",

			"UNT",
		}, true, maxGroupRepeatCountExceeded, nil,
	},
}

func TestSegSeqValidator1(t *testing.T) {
	msgSpec := getMsgSpec("AUTHOR_D.14B")

	for _, spec := range authorSegSeqSpec {
		fmt.Printf(">>>>>>>>>>>>>>>>>>> spec: %#v\n", spec)
		validator := NewSegSeqValidator(msgSpec)
		require.NotNil(t, validator)
		segments := mapToSegs(spec.segmentIDs)
		require.NotNil(t, segments)
		rawMessage := msg.NewRawMsg("AUTHOR", segments)
		fmt.Printf("Validating raw message: %s", rawMessage)
		nestedMsg, err := validator.Validate(rawMessage)

		if spec.expectError {
			assert.NotNil(t, err)
			assert.Nil(t, nestedMsg)
			fmt.Printf("Expected error kind: %s\n", spec.errorKind)
			fmt.Printf("\tError was: %s\n", err)
			err, ok := err.(SegSeqError)
			assert.True(t, ok)
			require.Equal(t, spec.errorKind, err.kind)
		} else {
			require.Nil(t, err)
			require.NotNil(t, nestedMsg)
			fmt.Printf("Constructed nested message:\n%s", nestedMsg.Dump())
			if spec.validateNestedMsg != nil {
				spec.validateNestedMsg(t, nestedMsg)
			}
		}
	}
}

func TestConsumeEmpty(t *testing.T) {
	defer func() {
		// Call to consume should panic
		if r := recover(); r != nil {
			fmt.Printf("recovered in TestConsumeEmpty\n")
		} else {
			t.Fail()
		}
	}()
	msgSpec := getMsgSpec("AUTHOR_D.14B")
	validator := NewSegSeqValidator(msgSpec)

	validator.segs = []*msg.Seg{}
	validator.consumeMulti()
}

var segABC *msg.Seg = msg.NewSeg("ABC")
var segDEF *msg.Seg = msg.NewSeg("DEF")

var consumeSpec = []struct {
	segsBefore []*msg.Seg
	segsAfter  []*msg.Seg
}{
	{
		[]*msg.Seg{segABC},
		[]*msg.Seg{},
	},
	{
		[]*msg.Seg{segABC, segABC},
		[]*msg.Seg{},
	},
	{
		[]*msg.Seg{segABC, segDEF},
		[]*msg.Seg{segDEF},
	},
	{
		[]*msg.Seg{segABC, segABC, segDEF},
		[]*msg.Seg{segDEF},
	},
	{
		[]*msg.Seg{segABC, segDEF, segABC},
		[]*msg.Seg{segDEF, segABC},
	},
}

func TestConsumeNonEmpty(t *testing.T) {
	msgSpec := getMsgSpec("AUTHOR_D.14B")

	for _, spec := range consumeSpec {
		validator := NewSegSeqValidator(msgSpec)
		validator.segs = spec.segsBefore
		validator.consumeMulti()
		assert.Equal(t, spec.segsAfter, validator.segs)
	}
}

func BenchmarkValidateSeq(b *testing.B) {
	segmentIDs := []string{
		"UNH", "BGM",
		"DTM", "BUS",
		// Group 4
		"LIN", "LIN", "LIN", "LIN",
		// Group 7
		"FII", "CTA", "COM", "COM", "COM",
		"FII", "CTA", "COM", "COM", "COM",

		"UNT",
	}
	msgSpec := getMsgSpec("AUTHOR_D.14B")
	validator := NewSegSeqValidator(msgSpec)
	require.NotNil(b, validator)
	segments := mapToSegs(segmentIDs)
	require.NotNil(b, segments)
	rawMessage := msg.NewRawMsg("AUTHOR", segments)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nestedMsg, err := validator.Validate(rawMessage)
		require.Nil(b, err)
		require.NotNil(b, nestedMsg)
	}
}

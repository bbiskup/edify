package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMsgTypeFromUNH(t *testing.T) {
	seg := rawmsg.NewRawSeg("UNH")
	seg.AddElem(rawmsg.NewRawDataElem([]string{"123"}))
	seg.AddElem(rawmsg.NewRawDataElem([]string{"ABC", "x", "y"}))
	msgType, err := getMsgTypeFromUNH(seg)
	assert.Nil(t, err)
	assert.Equal(t, "ABC", msgType)
}

func TestGetMsgTypeFromUNT(t *testing.T) {
	seg := rawmsg.NewRawSeg("UNT")
	seg.AddElem(rawmsg.NewRawDataElem([]string{"2"}))
	msgType, err := getSegCountFromUNT(seg)
	assert.Nil(t, err)
	assert.Equal(t, 2, msgType)
}

var validMsgTestSpecs = []struct {
	fileName string
	checkFn  func(*testing.T, *msg.NestedMsg)
}{
	// TODO support repeating data elems e.g. PAXLST_1.txt, 1st COM (M3)
	// COM+703-555-1212:TE+703-555-4545:FX'
	//{"PAXLST_1.txt"},
	//{"PAXLST_2.txt"},

	// TODO Implement other checkFns
	{"CUSRES_1.txt",
		func(t *testing.T, nestedMsg *msg.NestedMsg) {
			//assert.Equal(t, 4, nestedMsg.Count())
			fmt.Printf("nestedMsg: %s", nestedMsg.Dump())
			assert.Equal(t, 5, nestedMsg.GetTopLevelGrp().Count())
			assert.Equal(t, "UNH", nestedMsg.GetTopLevelGrp().GetPart(0).Id())

			group3Part := nestedMsg.GetTopLevelGrp().GetPart(2).(*msg.RepSegGrp)
			group3 := group3Part.GetSegGrp(0)
			assert.Equal(t, "Group_3", group3Part.Id())
			rff := group3.GetPart(0).(*msg.RepSeg)
			assert.Equal(t, "RFF", rff.Id())
			rff0 := rff.GetSeg(0)
			assert.Equal(t, "TN", rff0.DataElems[0].(*msg.CompositeDataElem).SimpleDataElems[0].Value)
		},
	},
	{"CUSRES_2.txt", nil},
	{"INVOIC_1.txt", nil},
	{"INVOIC_2.txt", nil},
	{"INVOIC_3.txt", nil},

	{"ORDERS_1.txt", nil},
	{"ORDERS_2.txt", nil},
	{"ORDERS_3.txt", nil},
}

func TestValidateMsg(t *testing.T) {
	validator := GetTestMsgValidator(t)
	assert.Equal(t, 194, validator.MsgSpecCount())
	assert.Equal(t, 156, validator.SegSpecCount())

	for _, testSpec := range validMsgTestSpecs {
		fileName := "../../testdata/messages/" + testSpec.fileName
		fmt.Printf("EDIFACT file: %s", fileName)
		rawMsg, err := GetRawMsg(fileName)
		require.Nil(t, err)
		nestedMsg, err := validator.Validate(rawMsg)
		assert.NotNil(t, nestedMsg)
		require.Nil(t, err)

		if testSpec.checkFn != nil {
			testSpec.checkFn(t, nestedMsg)
		}
	}
}

func BenchmarkValidateINVOICMsg(b *testing.B) {
	validator := GetTestMsgValidator(b)
	rawMsg, err := GetRawMsg("../../testdata/messages/INVOIC_1.txt")
	require.Nil(b, err)

	for i := 0; i < b.N; i++ {
		nestedMsg, err := validator.Validate(rawMsg)
		assert.NotNil(b, nestedMsg)
		require.Nil(b, err)
	}
}

func BenchmarkGetTestMsgValidator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetTestMsgValidator(b)
	}
}

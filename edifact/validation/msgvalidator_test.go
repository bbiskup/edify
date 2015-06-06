package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/spec/specparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getRawMsg(fileName string) (*msg.RawMsg, error) {
	parser := msg.NewParser()
	return parser.ParseRawMsgFile(fileName)
}

func TestGetMsgTypeFromUNH(t *testing.T) {
	seg := msg.NewSeg("UNH")
	seg.AddElem(msg.NewDataElem([]string{"123"}))
	seg.AddElem(msg.NewDataElem([]string{"ABC", "x", "y"}))
	msgType, err := getMsgTypeFromUNH(seg)
	assert.Nil(t, err)
	assert.Equal(t, "ABC", msgType)
}

func TestGetMsgTypeFromUNT(t *testing.T) {
	seg := msg.NewSeg("UNT")
	seg.AddElem(msg.NewDataElem([]string{"2"}))
	msgType, err := getSegCountFromUNT(seg)
	assert.Nil(t, err)
	assert.Equal(t, 2, msgType)
}

var validMsgTestSpecs = []struct {
	fileName string
}{
	// TODO support repeating data elems e.g. PAXLST_1.txt, 1st COM (M3)
	// COM+703-555-1212:TE+703-555-4545:FX'
	//{"PAXLST_1.txt"},
	//{"PAXLST_2.txt"},

	{"CUSRES_1.txt"},
	{"CUSRES_2.txt"},
	{"INVOIC_1.txt"}, // TODO errors (repetition of group 1 not detected correctly?)
	{"INVOIC_2.txt"},
	{"INVOIC_3.txt"},

	{"ORDERS_1.txt"},
	{"ORDERS_2.txt"},
	{"ORDERS_3.txt"},
}

func TestValidateMsg(t *testing.T) {
	parser, err := specparser.NewFullSpecParser("14B", "../../testdata/d14b")
	require.Nil(t, err)
	segSpecs, err := parser.ParseSegSpecsWithPrerequisites()
	require.Nil(t, err)

	msgSpecs, err := parser.ParseMsgSpecs(segSpecs)
	require.Nil(t, err)

	validator := NewMsgValidator(msgSpecs, segSpecs)

	for _, testSpec := range validMsgTestSpecs {
		rawMsg, err := getRawMsg("../../testdata/messages/" + testSpec.fileName)
		require.Nil(t, err)
		nestedMsg, err := validator.Validate(rawMsg)
		assert.NotNil(t, nestedMsg)
		require.Nil(t, err)
	}

}

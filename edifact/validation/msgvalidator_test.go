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
	// {"INVOIC_1.txt"},  // TODO errors (repetition of group 1 not detected correctly?)
	{"INVOIC_2.txt"},
	{"INVOIC_3.txt"},
	{"ORDERS_1.txt"},
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
		assert.Nil(t, err)
	}

}

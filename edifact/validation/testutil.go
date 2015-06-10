package validation

import (
	"github.com/bbiskup/edify/edifact/rawmsg"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	"github.com/bbiskup/edify/edifact/spec/specparser"
	"github.com/stretchr/testify/require"
	"testing"
)

func getMsgSpec(fileName string) *msp.MsgSpec {
	parser := msp.NewMsgSpecParser(&msp.MockSegSpecProviderImpl{})
	msgSpec, err := parser.ParseSpecFile("../../testdata/d14b/edmd/" + fileName)
	if err != nil {
		panic("spec is nil")
	}
	return msgSpec
}

func mapToRawSegs(segmentIDs []string) []*rawmsg.RawSeg {
	result := []*rawmsg.RawSeg{}
	for _, segmentID := range segmentIDs {
		result = append(result, rawmsg.NewRawSeg(segmentID))
	}
	return result
}

func GetRawMsg(fileName string) (*rawmsg.RawMsg, error) {
	parser := rawmsg.NewParser()
	return parser.ParseRawMsgFile(fileName)
}

func GetTestMsgValidator(tb testing.TB) *MsgValidator {
	parser, err := specparser.NewFullSpecParser("14B", "../../testdata/d14b")
	require.Nil(tb, err)
	segSpecs, err := parser.ParseSegSpecsWithPrerequisites()
	require.Nil(tb, err)

	msgSpecs, err := parser.ParseMsgSpecs(segSpecs)
	require.Nil(tb, err)

	return NewMsgValidator(msgSpecs, segSpecs)
}

package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/spec/message"
)

func getMsgSpec(fileName string) *message.MsgSpec {
	parser := message.NewMsgSpecParser(&message.MockSegSpecProviderImpl{})
	msgSpec, err := parser.ParseSpecFile("../../testdata/d14b/edmd/" + fileName)
	if err != nil {
		panic("spec is nil")
	}
	return msgSpec
}

func mapToSegs(segmentIDs []string) []*msg.Seg {
	result := []*msg.Seg{}
	for _, segmentID := range segmentIDs {
		result = append(result, msg.NewSeg(segmentID))
	}
	return result
}

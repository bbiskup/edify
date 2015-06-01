package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/spec/message"
)

func getMsgSpec(fileName string) *message.MsgSpec {
	parser := message.NewMsgSpecParser(&message.MockSegmentSpecProviderImpl{})
	msgSpec, err := parser.ParseSpecFile("../../testdata/d14b/edmd/" + fileName)
	if err != nil {
		panic("spec is nil")
	}
	return msgSpec
}

func mapToSegments(segmentIDs []string) []*msg.Segment {
	result := []*msg.Segment{}
	for _, segmentID := range segmentIDs {
		result = append(result, msg.NewSegment(segmentID))
	}
	return result
}

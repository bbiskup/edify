package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/spec/message"
)

func getMessageSpec(fileName string) *message.MessageSpec {
	parser := message.NewMessageSpecParser(&message.MockSegmentSpecProviderImpl{})
	messageSpec, err := parser.ParseSpecFile("../../testdata/d14b/edmd/" + fileName)
	if err != nil {
		panic("spec is nil")
	}
	return messageSpec
}

func mapToSegments(segmentIDs []string) []*msg.Segment {
	result := []*msg.Segment{}
	for _, segmentID := range segmentIDs {
		result = append(result, msg.NewSegment(segmentID))
	}
	return result
}

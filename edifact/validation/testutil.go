package validation

import (
	"github.com/bbiskup/edify/edifact/rawmsg"
	msp "github.com/bbiskup/edify/edifact/spec/message"
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

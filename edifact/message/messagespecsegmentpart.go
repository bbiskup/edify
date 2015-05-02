package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/segment"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment  specification in message specification
type MessageSpecSegmentPart struct {
	MessageSpecPartBase
	SegmentSpec *segment.SegmentSpec
}

func (p *MessageSpecSegmentPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment %d %s %s", p.MaxCount(), mandatoryStr, p.SegmentSpec.Name)
}

func (p *MessageSpecSegmentPart) IsGroup() bool {
	return false
}

func NewMessageSpecSegmentPart(segmentSpec *segment.SegmentSpec, maxCount int, isMandatory bool) *MessageSpecSegmentPart {
	return &MessageSpecSegmentPart{
		MessageSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
		},
		segmentSpec,
	}
}

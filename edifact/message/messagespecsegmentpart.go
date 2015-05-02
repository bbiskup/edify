package message

import (
	"github.com/bbiskup/edify/edifact/segment"
)

type MessageSpecSegmentPart struct {
	MessageSpecPartBase
	SegmentSpec *segment.SegmentSpec
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

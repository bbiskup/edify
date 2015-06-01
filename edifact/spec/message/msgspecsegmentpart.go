package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment  specification in message specification
type MsgSpecSegmentPart struct {
	MsgSpecPartBase
	SegmentSpec *segment.SegmentSpec
}

func (p *MsgSpecSegmentPart) Id() string {
	return p.SegmentSpec.Id
}

func (p *MsgSpecSegmentPart) Name() string {
	return p.SegmentSpec.Name
}

func (p *MsgSpecSegmentPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment %d %s %s", p.MaxCount(), mandatoryStr, p.SegmentSpec.Name)
}

func (p *MsgSpecSegmentPart) IsGroup() bool {
	return false
}

func NewMsgSpecSegmentPart(
	segmentSpec *segment.SegmentSpec,
	maxCount int, isMandatory bool, parent MsgSpecPart) *MsgSpecSegmentPart {

	return &MsgSpecSegmentPart{
		MsgSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		segmentSpec,
	}
}

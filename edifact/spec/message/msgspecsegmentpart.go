package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment  specification in message specification
type MsgSpecSegmentPart struct {
	MsgSpecPartBase
	SegSpec *segment.SegSpec
}

func (p *MsgSpecSegmentPart) Id() string {
	return p.SegSpec.Id
}

func (p *MsgSpecSegmentPart) Name() string {
	return p.SegSpec.Name
}

func (p *MsgSpecSegmentPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment %d %s %s", p.MaxCount(), mandatoryStr, p.SegSpec.Name)
}

func (p *MsgSpecSegmentPart) IsGroup() bool {
	return false
}

func NewMsgSpecSegmentPart(
	segSpec *segment.SegSpec,
	maxCount int, isMandatory bool, parent MsgSpecPart) *MsgSpecSegmentPart {

	return &MsgSpecSegmentPart{
		MsgSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		segSpec,
	}
}

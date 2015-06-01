package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment  specification in message specification
type MsgSpecSegPart struct {
	MsgSpecPartBase
	SegSpec *segment.SegSpec
}

func (p *MsgSpecSegPart) Id() string {
	return p.SegSpec.Id
}

func (p *MsgSpecSegPart) Name() string {
	return p.SegSpec.Name
}

func (p *MsgSpecSegPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment %d %s %s", p.MaxCount(), mandatoryStr, p.SegSpec.Name)
}

func (p *MsgSpecSegPart) IsGroup() bool {
	return false
}

func NewMsgSpecSegPart(
	segSpec *segment.SegSpec,
	maxCount int, isMandatory bool, parent MsgSpecPart) *MsgSpecSegPart {

	return &MsgSpecSegPart{
		MsgSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		segSpec,
	}
}

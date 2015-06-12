package message

import (
	"fmt"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/bbiskup/edify/edifact/util"
)

// Seg  specification in message specification
type MsgSpecSegPart struct {
	MsgSpecPartBase
	SegSpec *ssp.SegSpec
}

func NewMsgSpecSegPart(
	segSpec *ssp.SegSpec,
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

func (p *MsgSpecSegPart) Id() string {
	return p.SegSpec.Id
}

func (p *MsgSpecSegPart) Name() string {
	return p.SegSpec.Name
}

func (p *MsgSpecSegPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Seg %d %s %s", p.MaxCount(), mandatoryStr, p.SegSpec.Name)
}

func (p *MsgSpecSegPart) IsGroup() bool {
	return false
}

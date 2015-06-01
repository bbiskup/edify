package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment group specification in message specification
type MsgSpecSegmentGroupPart struct {
	MsgSpecPartBase
	name     string
	children []MsgSpecPart
}

func (p *MsgSpecSegmentGroupPart) Id() string {
	return p.TriggerSegmentPart().SegmentSpec.Id
}

func (p *MsgSpecSegmentGroupPart) Name() string {
	return p.name
}

func (p *MsgSpecSegmentGroupPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment group %s %d %s (%d children)", p.Name(), p.MaxCount(), mandatoryStr, p.Count())
}

func (p *MsgSpecSegmentGroupPart) IsGroup() bool {
	return true
}

func (p *MsgSpecSegmentGroupPart) Count() int {
	return len(p.children)
}

func (p *MsgSpecSegmentGroupPart) Children() []MsgSpecPart {
	return p.children
}

func (p *MsgSpecSegmentGroupPart) Append(msgSpecPart MsgSpecPart) {
	p.children = append(p.children, msgSpecPart)
}

// First segment spec contained in group. This is by definition
// a segment spec, not a new group.
func (p *MsgSpecSegmentGroupPart) TriggerSegmentPart() *MsgSpecSegmentPart {
	if len(p.children) > 0 {
		triggerSegmentPart, ok := p.children[0].(*MsgSpecSegmentPart)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", triggerSegmentPart))
		}
		return triggerSegmentPart
	} else {
		return nil
	}
}

func NewMsgSpecSegmentGroupPart(
	name string, children []MsgSpecPart,
	maxCount int, isMandatory bool, parent MsgSpecPart) *MsgSpecSegmentGroupPart {

	return &MsgSpecSegmentGroupPart{
		MsgSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		name,
		children,
	}
}

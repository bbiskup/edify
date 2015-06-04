package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment group specification in message specification
type MsgSpecSegGrpPart struct {
	MsgSpecPartBase
	name     string
	children []MsgSpecPart
}

func (p *MsgSpecSegGrpPart) Id() string {
	return p.children[0].Id()
}

func (p *MsgSpecSegGrpPart) Name() string {
	return p.name
}

func (p *MsgSpecSegGrpPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment group %s %d %s (%d children)", p.Name(), p.MaxCount(), mandatoryStr, p.Count())
}

func (p *MsgSpecSegGrpPart) IsGroup() bool {
	return true
}

func (p *MsgSpecSegGrpPart) Count() int {
	return len(p.children)
}

func (p *MsgSpecSegGrpPart) Children() []MsgSpecPart {
	return p.children
}

func (p *MsgSpecSegGrpPart) Append(msgSpecPart MsgSpecPart) {
	p.children = append(p.children, msgSpecPart)
}

// First segment spec contained in group. This is by definition
// a segment spec, not a new group.
func (p *MsgSpecSegGrpPart) TriggerSegPart() *MsgSpecSegPart {
	if len(p.children) > 0 {
		triggerSegPart, ok := p.children[0].(*MsgSpecSegPart)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", triggerSegPart))
		}
		return triggerSegPart
	} else {
		return nil
	}
}

func NewMsgSpecSegGrpPart(
	name string, children []MsgSpecPart,
	maxCount int, isMandatory bool, parent MsgSpecPart) *MsgSpecSegGrpPart {

	return &MsgSpecSegGrpPart{
		MsgSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		name,
		children,
	}
}

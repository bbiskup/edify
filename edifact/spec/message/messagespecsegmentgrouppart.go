package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

// Segment group specification in message specification
type MessageSpecSegmentGroupPart struct {
	MessageSpecPartBase
	name     string
	children []MessageSpecPart
}

func (p *MessageSpecSegmentGroupPart) Name() string {
	return p.name
}

func (p *MessageSpecSegmentGroupPart) String() string {
	mandatoryStr := util.CustBoolStr(p.IsMandatory(), "mand.", "cond.")
	return fmt.Sprintf("Segment group %s %d %s (%d children)", p.Name(), p.MaxCount(), mandatoryStr, p.Count())
}

func (p *MessageSpecSegmentGroupPart) IsGroup() bool {
	return true
}

func (p *MessageSpecSegmentGroupPart) Count() int {
	return len(p.children)
}

func (p *MessageSpecSegmentGroupPart) Children() []MessageSpecPart {
	return p.children
}

func (p *MessageSpecSegmentGroupPart) Append(messageSpecPart MessageSpecPart) {
	p.children = append(p.children, messageSpecPart)
}

func NewMessageSpecSegmentGroupPart(
	name string, children []MessageSpecPart,
	maxCount int, isMandatory bool, parent MessageSpecPart) *MessageSpecSegmentGroupPart {

	return &MessageSpecSegmentGroupPart{
		MessageSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
			parent:      parent,
		},
		name,
		children,
	}
}

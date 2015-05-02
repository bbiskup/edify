package message

type MessageSpecSegmentGroupPart struct {
	MessageSpecPartBase
	children []*MessageSpecPart
}

func (p *MessageSpecSegmentGroupPart) IsGroup() bool {
	return true
}

func (p *MessageSpecSegmentGroupPart) Count() int {
	return len(p.children)
}

func (p *MessageSpecSegmentGroupPart) Children() []*MessageSpecPart {
	return p.children
}

func (p *MessageSpecSegmentGroupPart) Append(messageSpecPart *MessageSpecPart) {
	p.children = append(p.children, messageSpecPart)
}

func NewMessageSpecSegmentGroupPart(children []*MessageSpecPart, maxCount int, isMandatory bool) *MessageSpecSegmentGroupPart {
	return &MessageSpecSegmentGroupPart{
		MessageSpecPartBase{
			maxCount:    maxCount,
			isMandatory: isMandatory,
		},
		children,
	}
}

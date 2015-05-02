package message

// Parts of a message specification (single segments or segment groups)
type MessageSpecPart interface {
	MaxCount() int
	IsMandatory() bool
	IsGroup() bool
}

// Parts of a message.
type MessageSpecPartBase struct {
	maxCount    int
	isMandatory bool
}

func (b *MessageSpecPartBase) MaxCount() int {
	return b.maxCount
}

func (b *MessageSpecPartBase) IsMandatory() bool {
	return b.isMandatory
}

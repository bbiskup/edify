package message

// Parts of a message specification (single segments or segment groups)
type MessageSpecPart interface {
	MinCount() int
	MaxCount() int
	IsMandatory() bool
	IsGroup() bool

	Id() string

	// returns nil when at nesting level 0
	Parent() MessageSpecPart
	String() string
	Name() string
}

// Parts of a message.
type MessageSpecPartBase struct {
	maxCount    int
	isMandatory bool
	parent      MessageSpecPart
}

func (b *MessageSpecPartBase) MinCount() int {
	if b.IsMandatory() {
		return 1
	} else {
		return 0
	}
}

func (b *MessageSpecPartBase) MaxCount() int {
	return b.maxCount
}

func (b *MessageSpecPartBase) IsMandatory() bool {
	return b.isMandatory
}

func (b *MessageSpecPartBase) Parent() MessageSpecPart {
	return b.parent
}

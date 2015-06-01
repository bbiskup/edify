package message

// Parts of a message specification (single segments or segment groups)
type MsgSpecPart interface {
	MinCount() int
	MaxCount() int
	IsMandatory() bool
	IsGroup() bool

	Id() string

	// returns nil when at nesting level 0
	Parent() MsgSpecPart
	String() string
	Name() string
}

// Parts of a message.
type MsgSpecPartBase struct {
	maxCount    int
	isMandatory bool
	parent      MsgSpecPart
}

func (b *MsgSpecPartBase) MinCount() int {
	if b.IsMandatory() {
		return 1
	} else {
		return 0
	}
}

func (b *MsgSpecPartBase) MaxCount() int {
	return b.maxCount
}

func (b *MsgSpecPartBase) IsMandatory() bool {
	return b.isMandatory
}

func (b *MsgSpecPartBase) Parent() MsgSpecPart {
	return b.parent
}

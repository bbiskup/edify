package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
)

// To be stored on stack while traversing group hierarchy
type SegSeqGroupContext struct {
	// may be nil (at top level)
	groupSpecPart *msgspec.MsgSpecSegGrpPart

	// top-level parts or children of current group
	parts []msgspec.MsgSpecPart

	// index on current level
	partIndex int

	// repeat count of current segment
	segmentRepeatCount int

	// repeat count of current group
	groupRepeatCount int

	// current group; nil if top-level
	repeatSegGroup *msg.RepSegGrp
}

func (c *SegSeqGroupContext) String() string {
	var groupName string
	if c.groupSpecPart != nil {
		groupName = c.groupSpecPart.Name()
	} else {
		groupName = "<toplevel>"
	}
	return fmt.Sprintf("currentGroupContext: %s (index: %d; # parts: %d)",
		groupName, c.partIndex, len(c.parts))
}

func (c *SegSeqGroupContext) AtEnd() bool {
	return c.partIndex >= len(c.parts)-1
}

func (c *SegSeqGroupContext) IsExhausted() bool {
	return c.partIndex >= len(c.parts)
}

func (c *SegSeqGroupContext) currentPart() msgspec.MsgSpecPart {
	return c.parts[c.partIndex]
}

func (c *SegSeqGroupContext) nextPart() msgspec.MsgSpecPart {
	return c.parts[c.partIndex+1]
}

func NewSegSeqGroupContext(
	groupSpecPart *msgspec.MsgSpecSegGrpPart,
	parts []msgspec.MsgSpecPart,
	repeatSegGroup *msg.RepSegGrp,
) *SegSeqGroupContext {

	return &SegSeqGroupContext{
		groupSpecPart, parts, 0, 0, 0, repeatSegGroup}
}

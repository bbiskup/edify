package validation

import (
	"fmt"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
)

// To be stored on stack while traversing group hierarchy
type SegSeqGroupContext struct {
	// may be nil (at top level)
	groupSpecPart *msgspec.MessageSpecSegmentGroupPart

	// top-level parts or chilren of current group
	parts []msgspec.MessageSpecPart

	// index on current level
	partIndex int

	// repeat count of current segment
	repeatCount int
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

func (c *SegSeqGroupContext) currentPart() msgspec.MessageSpecPart {
	return c.parts[c.partIndex]
}

func (c *SegSeqGroupContext) nextPart() msgspec.MessageSpecPart {
	return c.parts[c.partIndex+1]
}

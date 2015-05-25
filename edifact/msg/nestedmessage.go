package msg

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	noPartsText = "<no msg parts>"
)

// An EDIFACT message consisting of individual segments
// and possible nested segment groups.
// A nested message is suitable for element navigation
type NestedMessage struct {
	Name  string
	parts []SegmentOrGroup
}

func (m *NestedMessage) String() string {
	return fmt.Sprintf("NestedMessage %s (%d 1st-level parts)", m.Name, len(m.parts))
}

func (m *NestedMessage) segGroupDump(indent int, part SegmentOrGroup, buf *bytes.Buffer) {
	indentStr := strings.Repeat("\t", indent)
	switch part := part.(type) {
	case *Segment:
		buf.WriteString(fmt.Sprintf("%s%s\n", indentStr, part.Id()))
	case *SegmentGroup:
		buf.WriteString(fmt.Sprintf("%s%s\n", indentStr, part.Id()))
		for _, groupPart := range part.Parts {
			m.segGroupDump(indent+1, groupPart, buf)
		}
	default:
		panic(fmt.Sprintf("Unexpected type %T", part))
	}
}

// Comprehensive dump of segment/group structure
func (m *NestedMessage) SegGroupDump() string {
	if len(m.parts) == 0 {
		return noPartsText
	}
	var buf bytes.Buffer
	for _, part := range m.parts {
		m.segGroupDump(0, part, &buf)
	}
	return buf.String()
}

func NewNestedMessage(name string, parts []SegmentOrGroup) *NestedMessage {
	return &NestedMessage{name, parts}
}

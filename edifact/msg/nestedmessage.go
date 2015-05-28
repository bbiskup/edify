package msg

import (
	"bytes"
	"fmt"
)

const (
	noPartsText = "<no msg parts>"
)

// An EDIFACT message consisting of individual segments
// and possible nested segment groups.
// A nested message is suitable for element navigation
type NestedMessage struct {
	Name  string
	parts []RepeatMsgPart
}

func (m *NestedMessage) String() string {
	return fmt.Sprintf("NestedMessage %s (%d 1st-level parts)", m.Name, len(m.parts))
}

// Comprehensive dump of segment/group structure
func (m *NestedMessage) Dump(indent int) string {
	var buf bytes.Buffer
	if len(m.parts) == 0 {
		return noPartsText
	}
	buf.WriteString(fmt.Sprintf("%sMessage %s\n", getIndentStr(indent), m.Name))
	for _, part := range m.parts {
		buf.WriteString(part.Dump(indent + 1))
	}
	return buf.String()
}

func NewNestedMessage(name string, parts []RepeatMsgPart) *NestedMessage {
	return &NestedMessage{name, parts}
}

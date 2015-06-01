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
type NestedMsg struct {
	Name  string
	parts []RepeatMsgPart
}

func (m *NestedMsg) String() string {
	return fmt.Sprintf("NestedMsg %s (%d 1st-level parts)", m.Name, len(m.parts))
}

func (m *NestedMsg) Count() int {
	return len(m.parts)
}

func (m *NestedMsg) GetPart(index int) RepeatMsgPart {
	return m.parts[index]
}

func (m *NestedMsg) GetLastPart() RepeatMsgPart {
	return m.parts[len(m.parts)-1]
}

func (m *NestedMsg) AppendPart(part RepeatMsgPart) {
	m.parts = append(m.parts, part)
}

// Comprehensive dump of segment/group structure
func (m *NestedMsg) Dump() string {
	var buf bytes.Buffer
	if len(m.parts) == 0 {
		return noPartsText
	}
	buf.WriteString(fmt.Sprintf("Message %s\n", m.Name))
	for _, part := range m.parts {
		buf.WriteString(part.Dump(1))
	}
	return buf.String()
}

func NewNestedMsg(name string, parts []RepeatMsgPart) *NestedMsg {
	return &NestedMsg{name, parts}
}

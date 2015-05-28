package msg

import (
	"bytes"
	"fmt"
)

// A group of segments. Groups are not identified during message parsing,
// but during validation, when message specification is available
type SegmentGroup struct {
	id    string
	Parts []RepeatMsgPart
}

// From interface SegmentOrGroup
func (g *SegmentGroup) Id() string {
	return g.id
}

func NewSegmentGroup(id string, parts []RepeatMsgPart) *SegmentGroup {
	return &SegmentGroup{id, parts}
}

func (g *SegmentGroup) Dump(indent int) string {
	indentStr := getIndentStr(indent)
	var buf bytes.Buffer

	// Indentation of group name handled by parent RepeatSegmentGroup
	buf.WriteString(fmt.Sprintf("Group %s\n", g.Id()))
	for _, part := range g.Parts {
		buf.WriteString(fmt.Sprintf("%s  %s", indentStr, part.Dump(indent+1)))
	}
	return buf.String()
}

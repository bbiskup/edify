package msg

import (
	"bytes"
	"fmt"
)

// A group of segments. Groups are not identified during message parsing,
// but during validation, when message specification is available
type SegGrp struct {
	id    string
	parts []RepeatMsgPart
}

// From interface SegmentOrGroup
func (g *SegGrp) Id() string {
	return g.id
}

func (g *SegGrp) AppendSeg(segment *Segment) {
	g.parts = append(g.parts, NewRepeatSegment(segment))
}

func (g *SegGrp) AppendSegGroup(segGrp *RepeatSegGrp) {
	g.parts = append(g.parts, segGrp)
}

func NewSegGrp(id string, parts []RepeatMsgPart) *SegGrp {
	return &SegGrp{id, parts}
}

func (g *SegGrp) Dump(indent int) string {
	indentStr := getIndentStr(indent)
	var buf bytes.Buffer

	// Indentation of group name handled by parent RepeatSegGrp
	buf.WriteString(fmt.Sprintf("Group %s\n", g.Id()))
	for _, part := range g.parts {
		buf.WriteString(indentStr + "  " + part.Dump(indent+1))
	}
	return buf.String()
}

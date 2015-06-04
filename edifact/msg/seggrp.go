package msg

import (
	"bytes"
	"fmt"
)

// Map for looking up parts of a group by Id
type SegGrpMap map[string]RepeatMsgPart

// A group of segments. Groups are not identified during message parsing,
// but during validation, when message specification is available
type SegGrp struct {
	id      string
	parts   []RepeatMsgPart
	partMap SegGrpMap
}

// From interface SegOrGroup
func (g *SegGrp) Id() string {
	return g.id
}

func (g *SegGrp) Count() int {
	return len(g.parts)
}

func (g *SegGrp) GetPartByKey(key string) RepeatMsgPart {
	return g.partMap[key]
}

func (g *SegGrp) Contains(key string) bool {
	_, ok := g.partMap[key]
	return ok
}

// get n-th element
func (g *SegGrp) GetPart(index int) RepeatMsgPart {
	return g.parts[index]
}

func (g *SegGrp) AppendRepSeg(repSeg *RepSeg) {
	g.partMap[repSeg.Id()] = repSeg
	g.parts = append(g.parts, repSeg)
}

func (g *SegGrp) AppendRepSegGrp(repSegGrp *RepSegGrp) {
	g.partMap[repSegGrp.Id()] = repSegGrp
	g.parts = append(g.parts, repSegGrp)
}

func NewSegGrp(id string, parts ...RepeatMsgPart) *SegGrp {
	segGrpMap := make(SegGrpMap, len(parts))
	for _, part := range parts {
		segGrpMap[part.Id()] = part
	}
	return &SegGrp{id, parts, segGrpMap}
}

func (g *SegGrp) Dump(indent int) string {
	var buf bytes.Buffer

	// Indentation of group name handled by parent RepSegGrp
	buf.WriteString(fmt.Sprintf("Group %s\n", g.Id()))
	for _, part := range g.parts {
		buf.WriteString(part.Dump(indent + 1))
	}
	return buf.String()
}

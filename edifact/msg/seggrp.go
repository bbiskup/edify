package msg

import (
	"bytes"
	"errors"
	"fmt"
	msp "github.com/bbiskup/edify/edifact/spec/message"
)

// Map for looking up parts of a group by Id
type SegGrpMap map[string]RepMsgPart

// A group of segments. Groups are not identified during message parsing,
// but during validation, when message specification is available
type SegGrp struct {
	id      string
	parts   []RepMsgPart
	partMap SegGrpMap
}

func NewSegGrp(id string, parts ...RepMsgPart) *SegGrp {
	segGrpMap := make(SegGrpMap, len(parts))
	for _, part := range parts {
		segGrpMap[part.Id()] = part
	}
	return &SegGrp{id, parts, segGrpMap}
}

// From interface SegOrGroup
func (g *SegGrp) Id() string {
	return g.id
}

func (g *SegGrp) Count() int {
	return len(g.parts)
}

func (g *SegGrp) IsTopLevel() bool {
	return g.Id() == msp.TopLevelSegGroupName
}

func (g *SegGrp) GetPartByKey(key string) RepMsgPart {
	return g.partMap[key]
}

func (g *SegGrp) Contains(key string) bool {
	_, ok := g.partMap[key]
	return ok
}

// get n-th element
func (g *SegGrp) GetPart(index int) RepMsgPart {
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

func (g *SegGrp) FindNthOccurrenceOfSeg(segID string, occ int) (repSeg *RepSeg, err error) {
	for _, part := range g.parts {
		if part.Id() == segID {
			repSeg, ok := part.(*RepSeg)
			if !ok {
				return nil, errors.New(fmt.Sprintf(
					"Part %s is not a segment, but a %T", part.Id(), part))
			}
			return repSeg, nil
		}
	}
	return nil, errors.New(fmt.Sprintf(
		"Could not find segment %s in group %s", segID, g.Id()))
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

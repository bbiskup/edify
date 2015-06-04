package msg

import (
	"bytes"
	"strconv"
)

// A segment that is repeated 1 to n times
type RepSegGrp struct {
	groups []*SegGrp
}

// From Interface RepeatMsgPart
func (g *RepSegGrp) Count() int {
	return len(g.groups)
}

// From SegOrGroup
func (g *RepSegGrp) Id() string {
	return g.groups[0].Id()
}

// Append a repetition
func (g *RepSegGrp) Append(segGrp *SegGrp) {
	g.groups = append(g.groups, segGrp)
}

// Get n-th repetition
func (g *RepSegGrp) GetSegGrp(index int) *SegGrp {
	return g.groups[index]
}

// Get last repetition
func (g *RepSegGrp) GetLast() *SegGrp {
	return g.groups[len(g.groups)-1]
}

func (g *RepSegGrp) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	indentStr2 := getIndentStr(indent + 1)
	buf.WriteString(indentStr + "RepSegGrp\n")
	for repeat, group := range g.groups {
		buf.WriteString(
			indentStr2 + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
				group.Dump(indent+1))
	}
	return buf.String()
}

func NewRepSegGrp(groups ...*SegGrp) *RepSegGrp {
	return &RepSegGrp{
		groups,
	}
}

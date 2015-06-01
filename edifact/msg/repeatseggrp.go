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

func (g *RepSegGrp) Get(index int) *SegGrp {
	return g.groups[index]
}

func (g *RepSegGrp) GetLast() *SegGrp {
	return g.groups[len(g.groups)-1]
}

func (g *RepSegGrp) RepSegGrp(group *SegGrp) {
	g.groups = append(g.groups, group)
}

func (g *RepSegGrp) AppendSegGroupToLast(group *RepSegGrp) {
	lastGroup := g.groups[len(g.groups)-1]
	lastGroup.AppendSegGroup(group)
}

func (g *RepSegGrp) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	for repeat, group := range g.groups {
		buf.WriteString(
			indentStr + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
				group.Dump(indent))
	}
	return buf.String()
}

func NewRepSegGrp(groups ...*SegGrp) *RepSegGrp {
	return &RepSegGrp{
		groups,
	}
}

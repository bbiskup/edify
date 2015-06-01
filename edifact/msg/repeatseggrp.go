package msg

import (
	"bytes"
	"strconv"
)

// A segment that is repeated 1 to n times
type RepeatSegGrp struct {
	groups []*SegGrp
}

// From Interface RepeatMsgPart
func (g *RepeatSegGrp) Count() int {
	return len(g.groups)
}

// From SegOrGroup
func (g *RepeatSegGrp) Id() string {
	return g.groups[0].Id()
}

func (g *RepeatSegGrp) Get(index int) *SegGrp {
	return g.groups[index]
}

func (g *RepeatSegGrp) GetLast() *SegGrp {
	return g.groups[len(g.groups)-1]
}

func (g *RepeatSegGrp) RepeatSegGrp(group *SegGrp) {
	g.groups = append(g.groups, group)
}

func (g *RepeatSegGrp) AppendSegGroupToLast(group *RepeatSegGrp) {
	lastGroup := g.groups[len(g.groups)-1]
	lastGroup.AppendSegGroup(group)
}

func (g *RepeatSegGrp) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	for repeat, group := range g.groups {
		buf.WriteString(
			indentStr + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
				group.Dump(indent))
	}
	return buf.String()
}

func NewRepeatSegGrp(groups ...*SegGrp) *RepeatSegGrp {
	return &RepeatSegGrp{
		groups,
	}
}

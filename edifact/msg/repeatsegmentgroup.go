package msg

import (
	"bytes"
	"fmt"
)

// A segment that is repeated 1 to n times
type RepeatSegmentGroup struct {
	groups []*SegmentGroup
}

// From Interface RepeatMsgPart
func (g *RepeatSegmentGroup) Count() int {
	return len(g.groups)
}

// From SegmentOrGroup
func (g *RepeatSegmentGroup) Id() string {
	return g.groups[0].Id()
}

func (g *RepeatSegmentGroup) AddSegmentGroup(group *SegmentGroup) {
	g.groups = append(g.groups, group)
}

func (g *RepeatSegmentGroup) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	//buf.WriteString(fmt.Sprintf("%s%s\n", indentStr, g.groups[0].Id()))
	for repeat, group := range g.groups {
		buf.WriteString(fmt.Sprintf("%s[%d] %s\n", indentStr, repeat, group.Dump(indent)))
		//buf.WriteString(fmt.Sprintf("%s", group.Dump(indent+1)))
	}
	return buf.String()
}

func NewRepeatSegmentGroup(groups ...*SegmentGroup) *RepeatSegmentGroup {
	return &RepeatSegmentGroup{
		groups,
	}
}

package msg

import (
	"bytes"
	"fmt"
	"strconv"
)

// A segment that is repeated 1 to n times
type RepSegGrp struct {
	triggerSegId string
	groups       []*SegGrp
}

// From Interface RepeatMsgPart
func (g *RepSegGrp) Count() int {
	return len(g.groups)
}

// From SegOrGroup
func (g *RepSegGrp) Id() string {
	return g.triggerSegId
}

// Append a repetition
func (g *RepSegGrp) Append(segGrp *SegGrp) {
	checkGroupIdConsistency(g.triggerSegId, segGrp)
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

func checkGroupIdConsistency(triggerSegId string, groups ...*SegGrp) {
	for _, group := range groups {
		groupID := group.Id()
		if groupID != triggerSegId {
			panic(fmt.Sprintf("Inconsistent IDs: should be %s; got: %s",
				triggerSegId, groupID))
		}
	}
}

func NewRepSegGrp(triggerSegId string, groups ...*SegGrp) *RepSegGrp {
	checkGroupIdConsistency(triggerSegId, groups...)
	return &RepSegGrp{
		triggerSegId,
		groups,
	}
}

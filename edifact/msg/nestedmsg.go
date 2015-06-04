package msg

import (
	"bytes"
	"fmt"
)

const (
	TOP_LEVEL_GROUP = "_toplevel"
	noPartsText     = "<no msg parts>"
)

// An EDIFACT message consisting of individual segments
// and possible nested segment groups.
// A nested message is suitable for element navigation
type NestedMsg struct {
	Name string

	// The top level group does not actually exist in the specification,
	// but is used to allow uniform treatment of the parts of a nested
	// message at all levels.
	// By definition, the top level group is mandatory and
	// has a repeat count of 1
	TopLevelRepGrp *RepSegGrp
}

func (m *NestedMsg) String() string {
	return fmt.Sprintf("NestedMsg %s (%d 1st-level parts)",
		m.Name, m.Count())
}

// Number of elements at top level
func (m *NestedMsg) Count() int {
	return m.TopLevelRepGrp.GetSegGrp(0).Count()
}

func (m *NestedMsg) GetTopLevelGroup() *SegGrp {
	return m.TopLevelRepGrp.GetSegGrp(0)
}

// Comprehensive dump of segment/group structure
func (m *NestedMsg) Dump() string {
	var buf bytes.Buffer
	if m.Count() == 0 {
		return noPartsText
	}
	buf.WriteString(fmt.Sprintf(
		"Message %s\n%s", m.Name, m.TopLevelRepGrp.Dump(1)))
	return buf.String()
}

func NewNestedMsg(name string, parts ...RepeatMsgPart) *NestedMsg {
	topLevelGroup := NewSegGrp(TOP_LEVEL_GROUP, parts...)
	return &NestedMsg{
		Name:           name,
		TopLevelRepGrp: NewRepSegGrp(topLevelGroup)}
}

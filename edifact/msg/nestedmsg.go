package msg

import (
	"bytes"
	"fmt"
	msp "github.com/bbiskup/edify/edifact/spec/message"
)

const (
	noPartsText = "<no msg parts>\n"
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

func NewNestedMsg(name string, parts ...RepeatMsgPart) *NestedMsg {
	topLevelGroup := NewSegGrp(msp.TopLevelSegGroupName, parts...)
	return &NestedMsg{
		Name:           name,
		TopLevelRepGrp: NewRepSegGrp(topLevelGroup.Id(), topLevelGroup)}
}

func (m *NestedMsg) String() string {
	return fmt.Sprintf("NestedMsg %s (%d 1st-level parts)",
		m.Name, m.Count())
}

// Number of elements at top level
func (m *NestedMsg) Count() int {
	return m.TopLevelRepGrp.GetSegGrp(0).Count()
}

func (m *NestedMsg) GetTopLevelGrp() *SegGrp {
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
	buf.WriteString("\n")
	return buf.String()
}

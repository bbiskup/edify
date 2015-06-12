// Package rawmessage contains building blocks for unvalidated messages.
// Univalidated messages result from parsing purely by syntax, i.e. without
// using an EDIFACT message/segment/data element/code specification.
package rawmsg

import (
	"bytes"
	"fmt"
)

// A raw message contains a sequence of segments without
// the notion of segment groups, i.e. without nesting
type RawMsg struct {
	Name    string
	RawSegs []*RawSeg
}

func NewRawMsg(id string, rawSegs []*RawSeg) *RawMsg {
	return &RawMsg{id, rawSegs}
}

func (m *RawMsg) String() string {
	var buf bytes.Buffer
	for _, s := range m.RawSegs {
		buf.WriteString(fmt.Sprintf("\t%s", s.String()))
	}

	segmentsStr := buf.String()
	if len(segmentsStr) > 0 {
		return fmt.Sprintf("%s\n%s", m.Name, segmentsStr)
	} else {
		return m.Name
	}
}

func (m *RawMsg) AddRawSeg(rawSeg *RawSeg) {
	m.RawSegs = append(m.RawSegs, rawSeg)
}

// List of IDs of segments, in the order of their appearance.
// A separate entry is returned for each instance of a segment
func (m *RawMsg) RawSegIds() []string {
	result := []string{}
	for _, rawSeg := range m.RawSegs {
		result = append(result, rawSeg.Id())
	}
	return result
}

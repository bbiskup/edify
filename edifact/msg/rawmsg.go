package msg

import (
	"bytes"
	"fmt"
)

// A raw message contains a sequence of segments without
// the notion of segment groups, i.e. without nesting
type RawMessage struct {
	Name     string
	Segs []*Seg
}

func (m *RawMessage) String() string {
	var buf bytes.Buffer
	for _, s := range m.Segs {
		buf.WriteString(fmt.Sprintf("\t%s", s.String()))
	}

	segmentsStr := buf.String()
	if len(segmentsStr) > 0 {
		return fmt.Sprintf("%s\n%s", m.Name, segmentsStr)
	} else {
		return m.Name
	}
}

func (m *RawMessage) AddSeg(segment *Seg) {
	m.Segs = append(m.Segs, segment)
}

// List of IDs of segments, in the order of their appearance.
// A separate entry is returned for each instance of a segment
func (m *RawMessage) SegIds() []string {
	result := []string{}
	for _, segment := range m.Segs {
		result = append(result, segment.Id())
	}
	return result
}

func NewRawMessage(id string, segments []*Seg) *RawMessage {
	return &RawMessage{id, segments}
}

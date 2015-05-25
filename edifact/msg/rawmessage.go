package msg

// A raw message contains a sequence of segments without
// the notion of segment groups, i.e. without nesting
type RawMessage struct {
	Name     string
	Segments []*Segment
}

func (m *RawMessage) String() string {
	segmentsStr := ""
	for _, s := range m.Segments {
		segmentsStr += "\t" + s.String() + "\n"
	}
	return m.Name + "\n" + segmentsStr

}

func (m *RawMessage) AddSegment(segment *Segment) {
	m.Segments = append(m.Segments, segment)
}

// List of IDs of segments, in the order of their appearance.
// A separate entry is returned for each instance of a segment
func (m *RawMessage) SegmentIds() []string {
	result := []string{}
	for _, segment := range m.Segments {
		result = append(result, segment.Id())
	}
	return result
}

func NewRawMessage(id string, segments []*Segment) *RawMessage {
	return &RawMessage{id, segments}
}

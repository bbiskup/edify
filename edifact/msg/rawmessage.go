package msg

type Message struct {
	Name     string
	Segments []*Segment
}

func (m *Message) String() string {
	segmentsStr := ""
	for _, s := range m.Segments {
		segmentsStr += "\t" + s.String() + "\n"
	}
	return m.Name + "\n" + segmentsStr

}

func (m *Message) AddSegment(segment *Segment) {
	m.Segments = append(m.Segments, segment)
}

// List of IDs of segments, in the order of their appearance.
// A separate entry is returned for each instance of a segment
func (m *Message) SegmentIds() []string {
	result := []string{}
	for _, segment := range m.Segments {
		result = append(result, segment.Id())
	}
	return result
}

func NewMessage(id string, segments []*Segment) *Message {
	return &Message{id, segments}
}

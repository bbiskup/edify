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

func NewMessage(name string, segments []*Segment) *Message {
	return &Message{name, segments}
}

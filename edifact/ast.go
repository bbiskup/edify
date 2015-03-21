package edifact

// Abstract syntax tree for EDIFACT documents

type Element struct {
	Name  string
	Value string
}

func (e *Element) String() string {
	return e.Name + " " + e.Value
}

func NewElement(name string, value string) *Element {
	return &Element{name, value}
}

type Segment struct {
	Name     string
	Elements []*Element
}

func (s *Segment) String() string {
	elementsStr := ""
	for _, e := range s.Elements {
		elementsStr += "\t\t" + e.String() + "\n"
	}
	return s.Name + "\n" + elementsStr
}

func (s *Segment) AddElement(element *Element) {
	s.Elements = append(s.Elements, element)
}

func NewSegment(name string) *Segment {
	return &Segment{name, []*Element{}}
}

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

func NewMessage(name string) *Message {
	return &Message{name, []*Segment{}}
}

type Interchange struct {
	Messages []*Message
}

func (i *Interchange) String() string {
	result := ""
	for _, m := range i.Messages {
		result += "\n" + m.String() + "\n"
	}
	return result
}

func (i *Interchange) AddMessage(message *Message) {
	i.Messages = append(i.Messages, message)
}

func NewInterchange() *Interchange {
	return &Interchange{}
}

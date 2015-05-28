package msg

import (
	"bytes"
	"fmt"
)

type Segment struct {
	id       string
	Elements []*DataElement
}

// From interface SegmentOrGroup
func (g *Segment) Id() string {
	return g.id
}

func (s *Segment) String() string {
	var buf bytes.Buffer
	for _, e := range s.Elements {
		buf.WriteString(fmt.Sprintf("\t\t%s\n", e.String()))
	}
	return fmt.Sprintf("%s\n%s", s.id, buf.String())
}

func (s *Segment) AddElement(element *DataElement) {
	s.Elements = append(s.Elements, element)
}

func (s *Segment) AddElements(elements []*DataElement) {
	s.Elements = elements
}

func (s *Segment) Dump(indent int) string {
	indentStr := getIndentStr(indent)
	return fmt.Sprintf("%sSegment %s\n", indentStr, s.Id())
}

func NewSegment(id string) *Segment {
	return &Segment{id, []*DataElement{}}
}

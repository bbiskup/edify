package msg

import (
	"bytes"
	"fmt"
)

type Segment struct {
	id       string
	Elems []*DataElem
}

// From interface SegmentOrGroup
func (g *Segment) Id() string {
	return g.id
}

func (s *Segment) String() string {
	var buf bytes.Buffer
	for _, e := range s.Elems {
		buf.WriteString("\t\t" + e.String() + "\n")
	}
	return fmt.Sprintf("%s\n%s", s.id, buf.String())
}

func (s *Segment) AddElem(element *DataElem) {
	s.Elems = append(s.Elems, element)
}

func (s *Segment) AddElems(elements []*DataElem) {
	s.Elems = elements
}

func (s *Segment) Dump(indent int) string {
	indentStr := getIndentStr(indent)
	return fmt.Sprintf("%sSegment %s\n", indentStr, s.Id())
}

func NewSegment(id string) *Segment {
	return &Segment{id, []*DataElem{}}
}

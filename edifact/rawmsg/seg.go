package msg

import (
	"bytes"
	"fmt"
)

type Seg struct {
	id    string
	Elems []*RawDataElem
}

// From interface SegOrGroup
func (g *Seg) Id() string {
	return g.id
}

func (s *Seg) String() string {
	var buf bytes.Buffer
	for _, e := range s.Elems {
		buf.WriteString("\t\t" + e.String() + "\n")
	}
	return fmt.Sprintf("%s\n%s", s.id, buf.String())
}

func (s *Seg) AddElem(element *RawDataElem) {
	s.Elems = append(s.Elems, element)
}

func (s *Seg) AddElems(elements []*RawDataElem) {
	s.Elems = elements
}

func (s *Seg) Dump(indent int) string {
	indentStr := getIndentStr(indent)
	return fmt.Sprintf("%sSeg %s\n", indentStr, s.Id())
}

func NewSeg(id string) *Seg {
	return &Seg{id, []*RawDataElem{}}
}

package msg

import (
	"fmt"
)

// Composite data element after validation
type CompositeDataElem struct {
	id              string
	SimpleDataElems []*SimpleDataElem
}

func (e *CompositeDataElem) Id() string {
	return e.id
}

func (e *CompositeDataElem) String() string {
	return fmt.Sprintf("CompositeDataElem %s (%d simple data elems)", e.id, len(e.SimpleDataElems))
}

func NewCompositeDataElem(id string, simpleDataElems ...*SimpleDataElem) *CompositeDataElem {
	return &CompositeDataElem{
		id:              id,
		SimpleDataElems: simpleDataElems,
	}
}

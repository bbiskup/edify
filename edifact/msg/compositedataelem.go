package msg

import (
	"fmt"
)

// Composite data element after validation
type CompositeDataElem struct {
	id              string
	SimpleDataElems []*SimpleDataElem
}

func NewCompositeDataElem(id string, simpleDataElems ...*SimpleDataElem) *CompositeDataElem {
	return &CompositeDataElem{
		id:              id,
		SimpleDataElems: simpleDataElems,
	}
}

func (e *CompositeDataElem) Id() string {
	return e.id
}

func (s *CompositeDataElem) GetSimpleDataElemById(dataElemId string) (*SimpleDataElem, error) {
	for _, dataElem := range s.SimpleDataElems {
		if dataElem.Id() == dataElemId {
			return dataElem, nil
		}
	}
	return nil, fmt.Errorf("Data element '%s' not found in segment %s", dataElemId, s.id)
}

func (e *CompositeDataElem) String() string {
	return fmt.Sprintf("CompositeDataElem %s (%d simple data elems)", e.id, len(e.SimpleDataElems))
}

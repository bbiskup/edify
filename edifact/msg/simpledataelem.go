package msg

import (
	"fmt"
)

// A validated simple data element
type SimpleDataElem struct {
	id    string
	Value string
}

func (e *SimpleDataElem) Id() string {
	return e.id
}

func (e *SimpleDataElem) String() string {
	return fmt.Sprintf("SimpleDataElem %s: '%s'", e.Id(), e.Value)
}

func NewSimpleDataElem(id string, value string) *SimpleDataElem {
	return &SimpleDataElem{id: id, Value: value}
}

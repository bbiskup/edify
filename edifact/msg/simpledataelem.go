package msg

import (
	"fmt"
)

// A validated simple data element
type SimpleDataElem struct {
	id    string
	Value string
}

func NewSimpleDataElem(id string, value string) *SimpleDataElem {
	return &SimpleDataElem{id: id, Value: value}
}

func (e *SimpleDataElem) Id() string {
	return e.id
}

func (e *SimpleDataElem) String() string {
	return fmt.Sprintf("SimpleDataElem %s: '%s'", e.Id(), e.Value)
}

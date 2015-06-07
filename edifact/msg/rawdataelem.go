package msg

import (
	"fmt"
	"strings"
)

// Raw (unvalidated) data element.
// There is no associated spec information before validation
type RawDataElem struct {
	//Name   string
	Values []string
}

func (e *RawDataElem) IsSimple() bool {
	return len(e.Values) == 1
}

func (e *RawDataElem) buildComponentStr() string {
	result := []string{}
	for _, value := range e.Values {
		result = append(result, fmt.Sprintf("'%s'", value))
	}
	return strings.Join(result, " ")
}

func (e *RawDataElem) String() string {
	return fmt.Sprintf("RawDataElem %s", e.buildComponentStr())
}

func NewRawDataElem(values []string) *RawDataElem {
	return &RawDataElem{values}
}

package msg

import (
	"fmt"
	"strings"
)

type Element struct {
	//Name   string
	Values []string
}

func (e *Element) buildComponentStr() string {
	result := []string{}
	for _, value := range e.Values {
		result = append(result, fmt.Sprintf("'%s'", value))
	}
	return strings.Join(result, " ")
}

func (e *Element) String() string {
	//return e.Name + " " + strings.Join(e.Values, string(CompDataElemSep))
	return fmt.Sprintf("DataElement %s", e.buildComponentStr())
}

func NewElement(values []string) *Element {
	return &Element{values}
}

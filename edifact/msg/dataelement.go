package msg

import (
	"fmt"
	"strings"
)

type DataElement struct {
	//Name   string
	Values []string
}

func (e *DataElement) IsSimple() bool {
	return len(e.Values) == 1
}

func (e *DataElement) buildComponentStr() string {
	result := []string{}
	for _, value := range e.Values {
		result = append(result, fmt.Sprintf("'%s'", value))
	}
	return strings.Join(result, " ")
}

func (e *DataElement) String() string {
	//return e.Name + " " + strings.Join(e.Values, string(CompDataElemSep))
	return fmt.Sprintf("DataElement %s", e.buildComponentStr())
}

func NewDataElement(values []string) *DataElement {
	return &DataElement{values}
}

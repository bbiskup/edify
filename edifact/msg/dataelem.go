package msg

import (
	"fmt"
	"strings"
)

type DataElem struct {
	//Name   string
	Values []string
}

func (e *DataElem) IsSimple() bool {
	return len(e.Values) == 1
}

func (e *DataElem) buildComponentStr() string {
	result := []string{}
	for _, value := range e.Values {
		result = append(result, fmt.Sprintf("'%s'", value))
	}
	return strings.Join(result, " ")
}

func (e *DataElem) String() string {
	//return e.Name + " " + strings.Join(e.Values, string(CompDataElemSep))
	return fmt.Sprintf("DataElem %s", e.buildComponentStr())
}

func NewDataElem(values []string) *DataElem {
	return &DataElem{values}
}

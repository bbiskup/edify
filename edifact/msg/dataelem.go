package msg

import (
	"fmt"
	"strings"
)

// Used before a data element Id has been assigned during
// validation
const NoDataElemId = "_no_id"

type DataElem struct {
	//Name   string
	id     string
	Values []string
}

func (e *DataElem) Id() string {
	return e.id
}

func (e *DataElem) SetId(id string) {
	e.id = id
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
	return fmt.Sprintf("DataElem %s %s", e.id, e.buildComponentStr())
}

func NewDataElem(values []string) *DataElem {
	return &DataElem{NoDataElemId, values}
}

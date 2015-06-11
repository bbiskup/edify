package msg

import (
	"errors"
	"fmt"
)

// A segment after validation
type Seg struct {
	id        string
	DataElems []DataElem
}

func (s *Seg) Id() string {
	return s.id
}

func (s *Seg) AddDataElem(dataElem DataElem) {
	s.DataElems = append(s.DataElems, dataElem)
}

func (s *Seg) GetDataElemById(dataElemId string) (DataElem, error) {
	for _, dataElem := range s.DataElems {
		if dataElem.Id() == dataElemId {
			return dataElem, nil
		}
	}
	return nil, errors.New(fmt.Sprintf(
		"Data element '%s' not found in segment %s", dataElemId, s.id))
}

func (s *Seg) String() string {
	return fmt.Sprintf("Seg %s (%d data elems)", s.id, len(s.DataElems))
}

func NewSeg(id string, dataElems ...DataElem) *Seg {
	return &Seg{
		id:        id,
		DataElems: dataElems,
	}
}

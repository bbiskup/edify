package msg

import (
	"fmt"
)

// A segment after validation
type Seg struct {
	id        string
	DataElems []DataElem
}

func NewSeg(id string, dataElems ...DataElem) *Seg {
	return &Seg{
		id:        id,
		DataElems: dataElems,
	}
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
	return nil, fmt.Errorf("Data element '%s' not found in segment %s", dataElemId, s.id)
}

func (s *Seg) GetCompositeDataElemById(dataElemId string) (*CompositeDataElem, error) {
	dataElem, err := s.GetDataElemById(dataElemId)
	if err != nil {
		return nil, err
	}
	compositeDataElem, ok := dataElem.(*CompositeDataElem)
	if !ok {
		return nil, fmt.Errorf("Data element %s is not a composite data element", dataElemId)
	}
	return compositeDataElem, nil
}

func (s *Seg) GetSimpleDataElemById(dataElemId string) (*SimpleDataElem, error) {
	dataElem, err := s.GetDataElemById(dataElemId)
	if err != nil {
		return nil, err
	}
	simpleDataElem, ok := dataElem.(*SimpleDataElem)
	if !ok {
		return nil, fmt.Errorf("Data element %s is not a simple data element", dataElemId)
	}
	return simpleDataElem, nil
}

func (s *Seg) String() string {
	return fmt.Sprintf("Seg %s (%d data elems)", s.id, len(s.DataElems))
}

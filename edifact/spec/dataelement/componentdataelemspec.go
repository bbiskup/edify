package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

type ComponentDataElemSpec struct {
	Position    int
	IsMandatory bool

	// simple or composite data element spec
	SimpleDataElemSpec *SimpleDataElemSpec
}

func (s *ComponentDataElemSpec) String() string {
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf(
		"Component %s/%s @ %d (%s)",
		s.SimpleDataElemSpec.Id(), s.SimpleDataElemSpec.Name(),
		s.Position, isMandatoryStr)
}

func NewComponentDataElemSpec(position int, isMandatory bool, simpleDataElemSpec *SimpleDataElemSpec) *ComponentDataElemSpec {
	return &ComponentDataElemSpec{
		Position:           position,
		IsMandatory:        isMandatory,
		SimpleDataElemSpec: simpleDataElemSpec,
	}
}

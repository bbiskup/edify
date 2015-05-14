package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

type ComponentDataElementSpec struct {
	Position    int
	IsMandatory bool

	// simple or composite data element spec
	SimpleDataElemSpec *SimpleDataElementSpec
}

func (s *ComponentDataElementSpec) String() string {
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf(
		"Component %s/%s @ %d (%s)",
		s.SimpleDataElemSpec.Id(), s.SimpleDataElemSpec.Name(),
		s.Position, isMandatoryStr)
}

func NewComponentDataElementSpec(position int, isMandatory bool, simpleDataElementSpec *SimpleDataElementSpec) *ComponentDataElementSpec {
	return &ComponentDataElementSpec{
		Position:           position,
		IsMandatory:        isMandatory,
		SimpleDataElemSpec: simpleDataElementSpec,
	}
}

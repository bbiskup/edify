package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

type ComponentDataElementSpec struct {
	Position    int
	IsMandatory bool

	// simple or composite data element spec
	DataElemSpec DataElementSpec
}

func (s *ComponentDataElementSpec) String() string {
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf(
		"Component %s/%s @ %d (%s)",
		s.DataElemSpec.Id(), s.DataElemSpec.Name(),
		s.Position, isMandatoryStr)
}

func NewComponentDataElementSpec(position int, isMandatory bool, dataElementSpec DataElementSpec) *ComponentDataElementSpec {
	return &ComponentDataElementSpec{
		Position:     position,
		IsMandatory:  isMandatory,
		DataElemSpec: dataElementSpec,
	}
}

package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

type ComponentDataElementSpec struct {
	Num         int32
	IsMandatory bool
}

func (s *ComponentDataElementSpec) String() string {
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf("Component %d (%s)", s.Num, isMandatoryStr)
}

func NewComponentDataElementSpec(num int32, isMandatory bool) *ComponentDataElementSpec {
	return &ComponentDataElementSpec{
		Num:         num,
		IsMandatory: isMandatory,
	}
}

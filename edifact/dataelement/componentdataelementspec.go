package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
)

type ComponentDataElementSpec struct {
	Num         int32
	IsMandatory bool

	// simple or composite data element spec
	DataElementSpec interface{}
}

func (s *ComponentDataElementSpec) String() string {
	isMandatoryStr := util.CustBoolStr(s.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf("Component %d (%s)", s.Num, isMandatoryStr)
}

func NewComponentDataElementSpec(num int32, isMandatory bool, dataElementSpec interface{}) *ComponentDataElementSpec {
	return &ComponentDataElementSpec{
		Num:             num,
		IsMandatory:     isMandatory,
		DataElementSpec: dataElementSpec,
	}
}

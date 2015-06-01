package segment

import (
	"fmt"
	de "github.com/bbiskup/edify/edifact/spec/dataelement"
	"github.com/bbiskup/edify/edifact/util"
)

// A data element which is part of a segment specification
// (EDSD)
type SegmentDataElemSpec struct {
	DataElemSpec de.DataElemSpec
	Count        int
	IsMandatory  bool
}

func (e *SegmentDataElemSpec) String() string {
	mandatoryStr := util.CustBoolStr(e.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf("SegmentDataElem %s %dx %s", e.DataElemSpec.Id(), e.Count, mandatoryStr)
}

func NewSegmentDataElemSpec(
	dataElemSpec de.DataElemSpec, count int, isMandatory bool) *SegmentDataElemSpec {
	return &SegmentDataElemSpec{
		DataElemSpec: dataElemSpec,
		Count:        count,
		IsMandatory:  isMandatory,
	}
}

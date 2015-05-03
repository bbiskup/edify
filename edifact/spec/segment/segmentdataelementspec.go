package segment

import (
	"fmt"
	de "github.com/bbiskup/edify/edifact/spec/dataelement"
	"github.com/bbiskup/edify/edifact/util"
)

// A data element which is part of a segment specification
// (EDSD)
type SegmentDataElementSpec struct {
	DataElemSpec de.DataElementSpec
	Count        int
	IsMandatory  bool
}

func (e *SegmentDataElementSpec) String() string {
	mandatoryStr := util.CustBoolStr(e.IsMandatory, "mand.", "cond.")
	return fmt.Sprintf("SegmentDataElem %s %dx %s", e.DataElemSpec.Id(), e.Count, mandatoryStr)
}

func NewSegmentDataElementSpec(
	dataElemSpec de.DataElementSpec, count int, isMandatory bool) *SegmentDataElementSpec {
	return &SegmentDataElementSpec{
		DataElemSpec: dataElemSpec,
		Count:        count,
		IsMandatory:  isMandatory,
	}
}

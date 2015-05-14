package validation

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	spec_seg "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
)

type SegmentValidator struct {
	segmentSpecMap spec_seg.SegmentSpecMap
}

func (v *SegmentValidator) Validate(seg *msg.Segment) (isValid bool, err error) {
	spec := v.segmentSpecMap[seg.Id]
	if spec == nil {
		return false, errors.New(fmt.Sprintf("No spec for segment ID '%s'", seg.Id))
	}
	numDataElementSpecs := len(spec.SegmentDataElementSpecs)
	numDataElements := len(seg.Elements)
	if numDataElementSpecs != numDataElements {
		return false, errors.New(
			fmt.Sprintf("Incorrect number of data elements: got %d (%v), expected %d",
				numDataElements, seg.Elements, numDataElementSpecs))
	}
	for _, dataElemSpec := range spec.SegmentDataElementSpecs {
		log.Printf(" \t parsing data element %s", dataElemSpec.String())
	}
	return true, nil
}

func NewSegmentValidator(segmentSpecMap spec_seg.SegmentSpecMap) *SegmentValidator {
	return &SegmentValidator{segmentSpecMap: segmentSpecMap}
}

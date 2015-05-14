package validation

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	de "github.com/bbiskup/edify/edifact/spec/dataelement"
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

	return v.validateDataElems(
		spec.SegmentDataElementSpecs, seg.Elements)
}

func (v *SegmentValidator) validateDataElems(
	segmentDataElemSpecs []*spec_seg.SegmentDataElementSpec,
	dataElems []*msg.DataElement) (isVaid bool, err error) {

	for i, segDataElemSpec := range segmentDataElemSpecs {
		log.Printf(" \t parsing data element %s", segDataElemSpec.String())
		dataElem := dataElems[i]
		dataElemSpec := segDataElemSpec.DataElemSpec
		_, err := v.validateDataElem(dataElemSpec, dataElem)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// recursively validates data element spec
func (v *SegmentValidator) validateDataElem(
	dataElemSpec de.DataElementSpec, dataElem *msg.DataElement) (isValid bool, err error) {
	log.Printf("## dataElemSpec: %#v, dataElem: %v",
		dataElemSpec, dataElem)

	// TODO validate codes
	switch dataElemSpec := dataElemSpec.(type) {
	case *de.SimpleDataElementSpec:
		return dataElemSpec.Repr.Validate(dataElem.Values[0])
	case *de.CompositeDataElementSpec:
		for componentIndex, componentSpec := range dataElemSpec.ComponentSpecs {
			_, err := componentSpec.SimpleDataElemSpec.Repr.Validate(dataElem.Values[componentIndex])
			if err != nil {
				return false, err
			}
		}
		return true, nil
	default:
		panic("Invalid type")
	}
}

func NewSegmentValidator(segmentSpecMap spec_seg.SegmentSpecMap) *SegmentValidator {
	return &SegmentValidator{segmentSpecMap: segmentSpecMap}
}

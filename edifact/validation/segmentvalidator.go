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

func (v *SegmentValidator) Validate(seg *msg.Segment) error {
	spec := v.segmentSpecMap[seg.Id]
	if spec == nil {
		return errors.New(fmt.Sprintf("No spec for segment ID '%s'", seg.Id))
	}

	numDataElementSpecs := len(spec.SegmentDataElementSpecs)
	numDataElements := len(seg.Elements)
	if numDataElementSpecs != numDataElements {
		return errors.New(
			fmt.Sprintf("Incorrect number of data elements: got %d (%v), expected %d",
				numDataElements, seg.Elements, numDataElementSpecs))
	}

	return v.validateDataElems(
		spec.SegmentDataElementSpecs, seg.Elements)
}

func (v *SegmentValidator) validateDataElems(
	segmentDataElemSpecs []*spec_seg.SegmentDataElementSpec,
	dataElems []*msg.DataElement) error {

	for i, segDataElemSpec := range segmentDataElemSpecs {
		log.Printf(" \t parsing data element %s", segDataElemSpec.String())
		dataElem := dataElems[i]
		dataElemSpec := segDataElemSpec.DataElemSpec
		err := v.validateDataElem(dataElemSpec, dataElem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *SegmentValidator) validateSimpleDataElem(
	simpleDataElemSpec *de.SimpleDataElementSpec,
	value string) error {

	_, err := simpleDataElemSpec.Repr.Validate(value)
	if err != nil {
		return err
	}
	if simpleDataElemSpec.CodesSpecs != nil {
		if !simpleDataElemSpec.CodesSpecs.Contains(value) {
			return errors.New(
				fmt.Sprintf("Code %s not found", value))
		}
	}
	return nil
}

func (v *SegmentValidator) validateDataElem(
	dataElemSpec de.DataElementSpec, dataElem *msg.DataElement) error {
	log.Printf("## dataElemSpec: %#v, dataElem: %v",
		dataElemSpec, dataElem)

	// TODO validate codes
	switch dataElemSpec := dataElemSpec.(type) {
	case *de.SimpleDataElementSpec:
		return v.validateSimpleDataElem(dataElemSpec, dataElem.Values[0])
	case *de.CompositeDataElementSpec:
		for componentIndex, componentSpec := range dataElemSpec.ComponentSpecs {
			err := v.validateSimpleDataElem(
				componentSpec.SimpleDataElemSpec, dataElem.Values[componentIndex])
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic("Invalid type")
	}
}

func NewSegmentValidator(segmentSpecMap spec_seg.SegmentSpecMap) *SegmentValidator {
	return &SegmentValidator{segmentSpecMap: segmentSpecMap}
}

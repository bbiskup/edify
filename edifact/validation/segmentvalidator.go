package validation

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	de "github.com/bbiskup/edify/edifact/spec/dataelement"
	spec_seg "github.com/bbiskup/edify/edifact/spec/segment"
)

type SegmentValidator interface {
	Validate(seg *msg.Segment) error
}

// Validation of segments and their data elements
// The following aspects are validated:
// - cardinality of elements
// - correctness of representation (repr)
// - if a code mapping exists: validity of code
type SegmentValidatorImpl struct {
	segSpecMap spec_seg.SegSpecMap
}

func (v *SegmentValidatorImpl) Validate(seg *msg.Segment) error {
	spec := v.segSpecMap[seg.Id()]
	if spec == nil {
		return errors.New(fmt.Sprintf("No spec for segment ID '%s'", seg.Id()))
	}

	numDataElemSpecs := len(spec.SegmentDataElemSpecs)
	numDataElems := len(seg.Elements)
	if numDataElemSpecs != numDataElems {
		return errors.New(
			fmt.Sprintf("Incorrect number of data elements: got %d (%v), expected %d",
				numDataElems, seg.Elements, numDataElemSpecs))
	}

	return v.validateDataElems(
		spec.SegmentDataElemSpecs, seg.Elements)
}

func (v *SegmentValidatorImpl) validateDataElems(
	segmentDataElemSpecs []*spec_seg.SegmentDataElemSpec,
	dataElems []*msg.DataElem) error {

	for i, segDataElemSpec := range segmentDataElemSpecs {
		dataElem := dataElems[i]
		dataElemSpec := segDataElemSpec.DataElemSpec
		err := v.validateDataElem(dataElemSpec, dataElem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *SegmentValidatorImpl) validateSimpleDataElem(
	simpleDataElemSpec *de.SimpleDataElemSpec,
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

func (v *SegmentValidatorImpl) validateDataElem(
	dataElemSpec de.DataElemSpec, dataElem *msg.DataElem) error {

	// TODO validate codes
	switch dataElemSpec := dataElemSpec.(type) {
	case *de.SimpleDataElemSpec:
		return v.validateSimpleDataElem(dataElemSpec, dataElem.Values[0])
	case *de.CompositeDataElemSpec:
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

func NewSegmentValidatorImpl(segSpecMap spec_seg.SegSpecMap) *SegmentValidatorImpl {
	return &SegmentValidatorImpl{segSpecMap: segSpecMap}
}

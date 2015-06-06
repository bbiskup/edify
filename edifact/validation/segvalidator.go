package validation

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
)

type SegValidator interface {
	Validate(seg *msg.Seg) error
}

// Validation of segments and their data elements
// The following aspects are validated:
// - cardinality of elements
// - correctness of representation (repr)
// - if a code mapping exists: validity of code
type SegValidatorImpl struct {
	segSpecProvider ssp.SegSpecProvider
}

func (v *SegValidatorImpl) Validate(seg *msg.Seg) error {
	segID := seg.Id()
	log.Printf("Validating segment %s (%s)", segID, seg)
	spec := v.segSpecProvider.Get(segID)
	if spec == nil {
		return errors.New(fmt.Sprintf("No spec for segment ID '%s'", segID))
	}

	if ssp.IsUnValidatedSegment(segID) {
		log.Printf("######## Segment of type %s currently not validated", segID)
		return nil
	}

	minNumDataElemSpecs := spec.NumLeadingMandDataElems() // len(spec.SegDataElemSpecs)

	numDataElems := len(seg.Elems)
	if numDataElems < minNumDataElemSpecs {
		return errors.New(
			fmt.Sprintf("Seg %s: Incorrect number of data elements: got %d (%v), expected at least %d",
				segID, numDataElems, seg.Elems, minNumDataElemSpecs))
	}

	numDataElemSpecs := len(spec.SegDataElemSpecs)
	if numDataElems > numDataElemSpecs {
		return errors.New(fmt.Sprintf(
			"Too many data elements for segment %s: %d (should be only %d)",
			segID, numDataElems, numDataElemSpecs))
	}

	log.Printf("Validating %d top-level data elems: %s", numDataElems, seg.Elems)
	return v.validateDataElems(
		spec.SegDataElemSpecs, seg)
}

func (v *SegValidatorImpl) validateDataElems(
	segDataElemSpecs []*ssp.SegDataElemSpec,
	seg *msg.Seg) error {
	dataElems := seg.Elems

	for i, dataElem := range dataElems {
		segDataElemSpec := segDataElemSpecs[i]
		dataElemSpec := segDataElemSpec.DataElemSpec
		err := v.validateDataElem(segDataElemSpec, dataElemSpec, dataElem, segDataElemSpec.IsMandatory)
		if err != nil {
			return errors.New(fmt.Sprintf(
				"Error validating segment %s: %s", seg.Id(), err))
		}
	}
	return nil
}

func (v *SegValidatorImpl) validateSimpleDataElem(
	segDataElemSpec *ssp.SegDataElemSpec,
	simpleDataElemSpec *dsp.SimpleDataElemSpec,
	value string,
	isMandatory bool) error {

	log.Printf("Validating simple data elem %s", simpleDataElemSpec.Id())

	if value == "" {
		if isMandatory {
			return errors.New(fmt.Sprintf(
				"Missing value for mandatory simple data element %s", simpleDataElemSpec.Id()))
		} else {
			return nil
		}
	}

	_, err := simpleDataElemSpec.Repr.Validate(value)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Validation of repr %s for data element %s failed: %s",
			simpleDataElemSpec.Repr, segDataElemSpec.DataElemSpec.Id(), err))
	}
	if simpleDataElemSpec.CodesSpecs != nil {
		if !simpleDataElemSpec.CodesSpecs.Contains(value) {
			return errors.New(
				fmt.Sprintf("Code %s not found; code specs: %s",
					value, simpleDataElemSpec.CodesSpecs.CodeListStr()))
		}
	}
	return nil
}

func (v *SegValidatorImpl) validateDataElem(
	segDataElemSpec *ssp.SegDataElemSpec,
	dataElemSpec dsp.DataElemSpec,
	dataElem *msg.DataElem,
	isMandatory bool) error {
	log.Printf("Validating data elem %s", dataElem)
	log.Printf("\tSpec: %s", dataElemSpec)

	// TODO validate codes
	switch dataElemSpec := dataElemSpec.(type) {
	case *dsp.SimpleDataElemSpec:
		return v.validateSimpleDataElem(
			segDataElemSpec, dataElemSpec, dataElem.Values[0], isMandatory)
	case *dsp.CompositeDataElemSpec:
		log.Printf("\t%s is composite", dataElemSpec.Id())
		log.Printf("###### %d %s", len(dataElemSpec.ComponentSpecs), dataElemSpec.ComponentSpecs)
		for componentIndex, componentSpec := range dataElemSpec.ComponentSpecs {
			if componentIndex < len(dataElem.Values) {

				err := v.validateSimpleDataElem(
					segDataElemSpec,
					componentSpec.SimpleDataElemSpec,
					dataElem.Values[componentIndex],
					isMandatory && componentSpec.IsMandatory)
				if err != nil {
					return errors.New(fmt.Sprintf(
						"Error validating composite data elem spec %s (index %d): %s",
						dataElemSpec.Id(), componentIndex, err))
				}
			} else {
				break
			}

		}
		return nil
	default:
		panic("Invalid type")
	}
}

func NewSegValidatorImpl(segSpecProvider ssp.SegSpecProvider) *SegValidatorImpl {
	return &SegValidatorImpl{segSpecProvider: segSpecProvider}
}

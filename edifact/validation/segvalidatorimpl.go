package validation

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	// "log"
)

// Validation of segments and their data elements
// The following aspects are validated:
// - cardinality of elements
// - correctness of representation (repr)
// - if a code mapping exists: validity of code
type SegValidatorImpl struct {
	segSpecProvider ssp.SegSpecProvider
}

func NewSegValidatorImpl(segSpecProvider ssp.SegSpecProvider) *SegValidatorImpl {
	return &SegValidatorImpl{segSpecProvider: segSpecProvider}
}

// Checks data elements for correctness and constructs a (validated) Seg
func (v *SegValidatorImpl) Validate(rawSeg *rawmsg.RawSeg) (*msg.Seg, error) {
	segID := rawSeg.Id()
	// log.Printf("Validating segment %s (%s)", segID, seg)
	spec := v.segSpecProvider.Get(segID)
	if spec == nil {
		return nil, fmt.Errorf("No spec for segment ID '%s'", segID)
	}

	if ssp.IsUnValidatedSegment(segID) {
		//log.Printf("Segment of type %s currently not validated (validation not supported)", segID)
		return msg.NewSeg(segID), nil
	}

	minNumDataElemSpecs := spec.NumLeadingMandDataElems() // len(spec.SegDataElemSpecs)

	numDataElems := len(rawSeg.Elems)
	if numDataElems < minNumDataElemSpecs {
		return nil, fmt.Errorf(
			"Seg %s: Incorrect number of data elements: got %d (%v), expected at least %d",
			segID, numDataElems, rawSeg.Elems, minNumDataElemSpecs)
	}

	numDataElemSpecs := len(spec.SegDataElemSpecs)
	if numDataElems > numDataElemSpecs {
		return nil, fmt.Errorf(
			"Too many data elements for segment %s: %d (should be only %d)",
			segID, numDataElems, numDataElemSpecs)
	}

	seg := msg.NewSeg(rawSeg.Id())

	// log.Printf("Validating %d top-level data elems: %s", numDataElems, seg.Elems)
	err := v.validateDataElems(
		spec.SegDataElemSpecs, rawSeg, seg)

	if err != nil {
		return nil, err
	} else {
		return seg, nil
	}
}

func (v *SegValidatorImpl) validateDataElems(
	segDataElemSpecs []*ssp.SegDataElemSpec,
	rawSeg *rawmsg.RawSeg,
	seg *msg.Seg) error {
	dataElems := rawSeg.Elems

	for i, dataElem := range dataElems {
		segDataElemSpec := segDataElemSpecs[i]
		dataElemSpec := segDataElemSpec.DataElemSpec
		err := v.validateDataElem(
			segDataElemSpec, dataElemSpec,
			dataElem, segDataElemSpec.IsMandatory,
			seg)
		if err != nil {
			return fmt.Errorf("Error validating segment %s: %s", rawSeg.Id(), err)
		}
	}
	return nil
}

func (v *SegValidatorImpl) validateSimpleDataElem(
	segDataElemSpec *ssp.SegDataElemSpec,
	simpleDataElemSpec *dsp.SimpleDataElemSpec,
	value string,
	isMandatory bool) (simpleDataElem *msg.SimpleDataElem, err error) {

	// log.Printf("Validating simple data elem %s", simpleDataElemSpec.Id())

	if value == "" {
		if isMandatory {
			return nil, fmt.Errorf(
				"Missing value for mandatory simple data element %s",
				simpleDataElemSpec.Id())
		} else {
			// TODO: appropriate return value
			return nil, nil
		}
	}

	_, err = simpleDataElemSpec.Repr.Validate(value)
	if err != nil {
		return nil, fmt.Errorf(
			"Validation of repr %s for data element %s failed: %s",
			simpleDataElemSpec.Repr, segDataElemSpec.DataElemSpec.Id(), err)
	}
	if simpleDataElemSpec.CodesSpecs != nil {
		if !simpleDataElemSpec.CodesSpecs.Contains(value) {
			return nil, errors.New(
				fmt.Sprintf("Code %s not found; code specs: %s",
					value, simpleDataElemSpec.CodesSpecs.CodeListStr()))
		}
	}
	simpleDataElem = msg.NewSimpleDataElem(
		simpleDataElemSpec.Id(), value)

	return simpleDataElem, nil
}

func (v *SegValidatorImpl) validateDataElem(
	segDataElemSpec *ssp.SegDataElemSpec,
	dataElemSpec dsp.DataElemSpec,
	dataElem *rawmsg.RawDataElem,
	isMandatory bool,
	seg *msg.Seg) error {
	// log.Printf("Validating data elem %s", dataElem)
	// log.Printf("\tSpec: %s", dataElemSpec)

	switch dataElemSpec := dataElemSpec.(type) {
	case *dsp.SimpleDataElemSpec:
		simpleDataElem, err := v.validateSimpleDataElem(
			segDataElemSpec, dataElemSpec, dataElem.Values[0], isMandatory)
		if err == nil {
			seg.AddDataElem(simpleDataElem)
		}
		return err
	case *dsp.CompositeDataElemSpec:
		simpleDataElems := []*msg.SimpleDataElem{}
		// log.Printf("\t%s is composite", dataElemSpec.Id())
		// log.Printf("###### %d %s", len(dataElemSpec.ComponentSpecs), dataElemSpec.ComponentSpecs)
		for componentIndex, componentSpec := range dataElemSpec.ComponentSpecs {
			if componentIndex < len(dataElem.Values) {

				simpleDataElem, err := v.validateSimpleDataElem(
					segDataElemSpec,
					componentSpec.SimpleDataElemSpec,
					dataElem.Values[componentIndex],
					isMandatory && componentSpec.IsMandatory)
				if err == nil {
					simpleDataElems = append(simpleDataElems, simpleDataElem)
				} else {
					return fmt.Errorf(
						"Error validating composite data elem spec %s (index %d): %s",
						dataElemSpec.Id(), componentIndex, err)
				}
			} else {
				break
			}
		}
		compositeDataElem := msg.NewCompositeDataElem(dataElemSpec.Id(), simpleDataElems...)
		seg.AddDataElem(compositeDataElem)
		return nil
	default:
		panic("Invalid type")
	}
}

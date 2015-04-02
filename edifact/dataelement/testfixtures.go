package dataelement

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/codes"
)

func fixtureTextCodesSpecMap() codes.CodesSpecMap {
	codeSpec := codes.NewCodeSpec("id1", "codename", "codedescr")
	codesSpec := codes.NewCodesSpec("1000", "codesname", "codesdescr", []*codes.CodeSpec{codeSpec})
	specMap := codes.CodesSpecMap{"1000": codesSpec}
	return specMap
}

func fixtureSimpleDataElemSpec(id string) *SimpleDataElementSpec {
	result, err := NewSimpleDataElementSpec(
		id, fmt.Sprintf("name_%s", id), fmt.Sprintf("descr_%s", id),
		NewRepr(AlphaNum, true, 3), nil)
	if err != nil {
		panic(fmt.Sprintf("Fixture error: %s", err))
	}
	return result
}

func fixtureSimpleDataElems(ids ...string) SimpleDataElementSpecMap {
	result := SimpleDataElementSpecMap{}
	for _, id := range ids {
		result[id] = fixtureSimpleDataElemSpec(id)
	}
	return result
}

// Simple data element specs required for composite data element C001 (TRANSPORT MEANS)
// Code specs are omitted, so this fixture is suitable only for
// parsing, not validation.
// Non-key attributes use dummy values
func fixtureMultiSimpleDataElementSpecs() SimpleDataElementSpecMap {
	return fixtureSimpleDataElems("8179", "1131", "3055", "8178")
}

func fixtureCompositeDataElemSpec(id string) *CompositeDataElementSpec {
	return NewCompositeDataElementSpec(id, fmt.Sprintf("name_%s", id), fmt.Sprintf("descr_%s", id), nil)
}

func fixtureCompositeDataElems(ids ...string) CompositeDataElementSpecMap {
	result := CompositeDataElementSpecMap{}
	for _, id := range ids {
		result[id] = fixtureCompositeDataElemSpec(id)
	}
	return result
}

// Composite data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureCompositeDataElemsForSegmentADR() CompositeDataElementSpecMap {
	return fixtureCompositeDataElems("C817", "C090", "C819", "C517")
}

// Simple data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureSimpleDataElemsForSegmentADR() SimpleDataElementSpecMap {
	return fixtureSimpleDataElems("3164", "3251", "3207")
}

package dataelement

import (
	"fmt"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
)

func fixtureTextCodesSpecMap() csp.CodesSpecMap {
	codeSpec := csp.NewCodeSpec("id1", "codename", "codedescr")
	codesSpec := csp.NewCodesSpec("1000", "codesname", "codesdescr", []*csp.CodeSpec{codeSpec})
	specMap := csp.CodesSpecMap{"1000": codesSpec}
	return specMap
}

func fixtureSimpleDataElemSpec(id string) *SimpleDataElemSpec {
	result, err := NewSimpleDataElemSpec(
		id, fmt.Sprintf("name_%s", id), fmt.Sprintf("descr_%s", id),
		NewRepr(AlphaNum, true, 3), nil)
	if err != nil {
		panic(fmt.Sprintf("Fixture error: %s", err))
	}
	return result
}

func fixtureSimpleDataElems(ids ...string) SimpleDataElemSpecMap {
	result := SimpleDataElemSpecMap{}
	for _, id := range ids {
		result[id] = fixtureSimpleDataElemSpec(id)
	}
	return result
}

// Simple data element specs required for composite data element C001 (TRANSPORT MEANS)
// Code specs are omitted, so this fixture is suitable only for
// parsing, not validation.
// Non-key attributes use dummy values
func fixtureMultiSimpleDataElemSpecs() SimpleDataElemSpecMap {
	return fixtureSimpleDataElems("8179", "1131", "3055", "8178")
}

func fixtureCompositeDataElemSpec(id string) *CompositeDataElemSpec {
	return NewCompositeDataElemSpec(id, fmt.Sprintf("name_%s", id), fmt.Sprintf("descr_%s", id), nil)
}

func fixtureCompositeDataElems(ids ...string) CompositeDataElemSpecMap {
	result := CompositeDataElemSpecMap{}
	for _, id := range ids {
		result[id] = fixtureCompositeDataElemSpec(id)
	}
	return result
}

// Composite data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureCompositeDataElemsForSegADR() CompositeDataElemSpecMap {
	return fixtureCompositeDataElems("C817", "C090", "C819", "C517")
}

// Simple data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureSimpleDataElemsForSegADR() SimpleDataElemSpecMap {
	return fixtureSimpleDataElems("3164", "3251", "3207")
}

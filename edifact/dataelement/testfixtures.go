package dataelement

// Simple data element specs required for composite data element C001 (TRANSPORT MEANS)
// Code specs are omitted, so this fixture is suitable only for
// parsing, not validation.
// Non-key attributes use dummy values
func fixtureMultiSimpleDataElementSpecs() SimpleDataElementSpecMap {
	s_8179, _ := NewSimpleDataElementSpec("8179", "name_1", "descr_1", NewRepr(AlphaNum, true, 3), nil)
	s_1131, _ := NewSimpleDataElementSpec("1131", "name_2", "descr_2", NewRepr(AlphaNum, true, 3), nil)
	s_3055, _ := NewSimpleDataElementSpec("3055", "name_3", "descr_3", NewRepr(AlphaNum, true, 3), nil)
	s_8178, _ := NewSimpleDataElementSpec("8178", "name_4", "descr_4", NewRepr(AlphaNum, true, 3), nil)
	return SimpleDataElementSpecMap{
		"8179": s_8179,
		"1131": s_1131,
		"3055": s_3055,
		"8178": s_8178,
	}
}

/*
// Composite data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureCompositeDataElemsForSegmentADR() CompositeDataElementSpecMap {
	c_C817 := NewCompositeDataElementSpec("C817", "name_1", "descr_1", nil)
	c_C090 := NewCompositeDataElementSpec("C090", "name_1", "descr_1", nil)

	c_C817 := NewCompositeDataElementSpec("C817", "name_1", "descr_1", nil)
}

// Simple data element specs required for segmentspec ADR (ADDRESS)
// Component data elements are omitted, so this fixture is suitable only for
// parsing, not validation
// Non-key attributes use dummy values
func fixtureSimpleDataElemsForSegmentADR() SimpleDataElementSpecMap {
	s_3164 := NewCompositeDataElementSpec("3164", "name_1", "descr_1", nil)
	return SimpleDataElementSpecMap{
		3164: s_3164}
}*/

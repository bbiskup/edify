package dataelement

func fixtureMultiSimpleDataElementSpecs() SimpleDataElementSpecMap {
	s_8179, _ := NewSimpleDataElementSpec(8179, "name_1", "descr_1", NewRepr(AlphaNum, true, 3), nil)
	s_1131, _ := NewSimpleDataElementSpec(1131, "name_1", "descr_1", NewRepr(AlphaNum, true, 3), nil)
	s_3055, _ := NewSimpleDataElementSpec(3055, "name_1", "descr_1", NewRepr(AlphaNum, true, 3), nil)
	s_8178, _ := NewSimpleDataElementSpec(8178, "name_1", "descr_1", NewRepr(AlphaNum, true, 3), nil)
	return SimpleDataElementSpecMap{
		8179: s_8179,
		1131: s_1131,
		3055: s_3055,
		8178: s_8178,
	}
}

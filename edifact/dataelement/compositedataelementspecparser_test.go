package dataelement

import (
	"fmt"
	"testing"
)

var parserSpec = []struct {
	specLines      []string
	expectedResStr string
	expectErr      bool
}{
	{
		// Valid
		[]string{
			"       C001 TRANSPORT MEANS",
			"",
			"       Desc: Code and/or name identifying the type of means of",
			"             transport.",
			"",
			"010    8179  Transport means description code          C      an..8",
			"020    1131  Code list identification code             C      an..17",
			"030    3055  Code list responsible agency code         C      an..3",
			"040    8178  Transport means description               C      an..17",
		},
		"Composite C001 TRANSPORT MEANS 'Code an...'\n\tComponent 1131/name_1131 @ 20 (cond.)\n\tComponent 3055/name_3055 @ 30 (cond.)\n\tComponent 8178/name_8178 @ 40 (cond.)",
		false,
	},
	{
		// Invalid (no components)
		[]string{
			"       C001 TRANSPORT MEANS",
			"",
			"       Desc: Code and/or name identifying the type of means of",
			"             transport.",
			"",
		},
		"",
		true,
	},
	{
		// Invalid (no header)
		[]string{
			"",
			"       Desc: Code and/or name identifying the type of means of",
			"             transport.",
			"",
			"010    8179  Transport means description code          C      an..8",
			"020    1131  Code list identification code             C      an..17",
			"030    3055  Code list responsible agency code         C      an..3",
			"040    8178  Transport means description               C      an..17",
		},
		"",
		true,
	},
	{
		// Valid (fallback description)
		[]string{
			"       C001 TRANSPORT MEANS",
			"",
			"",
			"010    8179  Transport means description code          C      an..8",
			"020    1131  Code list identification code             C      an..17",
			"030    3055  Code list responsible agency code         C      an..3",
			"040    8178  Transport means description               C      an..17",
		},
		"Composite C001 TRANSPORT MEANS '<no des...'\n\tComponent 1131/name_1131 @ 20 (cond.)\n\tComponent 3055/name_3055 @ 30 (cond.)\n\tComponent 8178/name_8178 @ 40 (cond.)",
		true,
	},
}

func TestParser(t *testing.T) {
	for _, spec := range parserSpec {
		parser := NewCompositeDataElementSpecParser(fixtureMultiSimpleDataElementSpecs())
		res, err := parser.ParseSpec(spec.specLines)
		if err != nil && spec.expectErr {
			// fmt.Printf("expected err: %s", err)
			continue
		}

		if err != nil {
			t.Errorf(fmt.Sprintf("Failed to parse spec %s: %s", spec.specLines, err))
			continue
		}

		if res == nil {
			t.Errorf("No result")
			continue
		}

		resStr := res.String()
		if resStr != spec.expectedResStr {
			t.Errorf("Expected: %s, got: %s", spec.expectedResStr, resStr)
		}
	}
}

func TestParseFile(t *testing.T) {
	// TODO provide full data elements fixture
	parser := NewCompositeDataElementSpecParser(fixtureMultiSimpleDataElementSpecs())
	res, err := parser.ParseSpecFile("../../testdata/EDCD.14B_short")
	if err != nil {
		t.Errorf("Unable to parse composite data element spec: %s", err)
	}
	fmt.Printf("res: %s", res)
	if len(res) != 1 {
		t.Error("Expected 1 composite data element spec")
	}

	if res["C001"] == nil {
		t.Error("Composite data element C001 not found")
	}
}

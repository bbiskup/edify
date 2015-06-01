package dataelement

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedValid1 = `Composite C001 TRANSPORT MEANS 'Code an...'
	Component 1131/name_1131 @ 20 (cond.)
	Component 3055/name_3055 @ 30 (cond.)
	Component 8178/name_8178 @ 40 (cond.)`

const expectedValidFallbackDescription = `Composite C001 TRANSPORT MEANS '<no des...'
	Component 1131/name_1131 @ 20 (cond.)
	Component 3055/name_3055 @ 30 (cond.)
	Component 8178/name_8178 @ 40 (cond.)`

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
		expectedValid1,
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
		expectedValidFallbackDescription,
		true,
	},
}

func TestParser(t *testing.T) {
	for _, spec := range parserSpec {
		parser := NewCompositeDataElemSpecParser(fixtureMultiSimpleDataElemSpecs())
		res, err := parser.ParseSpec(spec.specLines)
		if err != nil && spec.expectErr {
			// fmt.Printf("expected err: %s", err)
			continue
		}

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, spec.expectedResStr, res.String())
	}
}

func TestParseFile(t *testing.T) {
	// TODO provide full data elements fixture
	parser := NewCompositeDataElemSpecParser(fixtureMultiSimpleDataElemSpecs())
	res, err := parser.ParseSpecFile("../../../testdata/EDCD.14B_short")
	assert.Nil(t, err)
	// fmt.Printf("res: %s", res)
	assert.Equal(t, 1, len(res))
	assert.NotNil(t, res["C001"])
}

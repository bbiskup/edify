package dataelement

import (
	"fmt"
	"testing"
)

/* OBSOLETE (old format using segment specification)
var parserSpec = []struct {
	specLines      []string
	expectedResStr string
	expectErr      bool
}{
	{
		// Valid
		[]string{
			"020    C138 PRICE MULTIPLIER INFORMATION               C    1      ",
			"       5394  Price multiplier rate                     M      n..12",
			"       5393  Price multiplier type code qualifier      C      an..3",
		},
		"Composite C138 PRICE MULTIPLIER INFORMATION 1 (conditional)\n\tComponent 5394 (mandatory)\n\tComponent 5393 (conditional)",
		false,
	},
	{
		// Valid, other values
		[]string{
			"021    C139 PRICE MULTIPLIER INFORMATION               M    2      ",
			"       5395  Price multiplier rate                     M      n..12",
			"       5396  Price multiplier type code qualifier      C      an..3",
			"",
		},
		"Composite C139 PRICE MULTIPLIER INFORMATION 2 (mandatory)\n\tComponent 5395 (mandatory)\n\tComponent 5396 (conditional)",
		false,
	},
	{
		// Invalid (incorrect header)
		[]string{
			"020    C138 PRICE MULTIPLIER INFORMATION               X    2      ",
			"       5394  Price multiplier rate                     M      n..12",
			"",
		},
		"",
		true,
	},
	{
		// Invalid (incorrect component)
		[]string{
			"020    C138 PRICE MULTIPLIER INFORMATION               C    2      ",
			"       5394  Price multiplier rate                     X      n..12",
			"",
		},
		"",
		true,
	},
	{
		// Invalid (no components)
		[]string{
			"020    C138 PRICE MULTIPLIER INFORMATION               C    2      ",
		},
		"",
		true,
	},
}

func TestParser(t *testing.T) {
	for _, spec := range parserSpec {
		parser := NewCompositeDataElementSpecParser()
		res, err := parser.ParseSpec(spec.specLines)
		if err != nil && spec.expectErr {
			fmt.Printf("expected err: %s", err)
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
*/

func TestParseFile(t *testing.T) {
	parser := NewCompositeDataElementSpecParser()
	res, err := parser.ParseSpecFile("../../testdata/EDCD.14B")
	if err != nil {
		t.Fatalf("Unable to parse composite data element spec: %s", err)
	}
	fmt.Printf("res: %s", res)
}

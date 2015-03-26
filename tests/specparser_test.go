package tests

import (
	sp "edifice/edifact/specparser"
	"fmt"
	"strings"
	"testing"
)

func TestSpecParser(t *testing.T) {
	p := sp.NewDataElementSpecParser()
	_, err := p.ParseSpecFile("../testdata/EDED.14B_short")
	if err != nil {
		t.Fatalf("Parse error: %s", err)
	}
}

const specLines = `
     1000  Document name                                           [B]

     Desc: Name of a document.

     Repr: an..35
`

func TestParseSpecLines(t *testing.T) {
	p := sp.NewDataElementSpecParser()
	res, err := p.ParseSpec(strings.Split(specLines, "\n"))
	if err != nil {
		t.Fatalf("Failed to parse specLines", err)
	}
	fmt.Printf("res: %s", res)
}

package tests

import (
	sp "edifice/edifact/specparser"
	"log"
	"reflect"
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
	log.Printf("res: %s", res)
}

func BenchmarkParseSpecLines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := sp.NewDataElementSpecParser()
		specs, err := p.ParseSpecFile("../testdata/EDED.14B")
		if err != nil {
			log.Printf("Parse error: %s\n", err)
			return
		}
		log.Printf("Parsed %d specs\n", len(specs))
	}
}

func TestSpecScannerFromFile(t *testing.T) {
	scanner, err := sp.NewSpecScanner("../testdata/specscanner/1")
	if err != nil {
		t.Fatalf("Error creating SpecScanner: %s", err)
	}

	expected := [][]string{
		[]string{"one"},
		[]string{"two"},
		[]string{"three"},
	}

	first, err := scanner.GetAllSpecLines()
	if err != nil {
		t.Fatalf("Error reading first spec line: %s", err)
	}
	if !reflect.DeepEqual(first, expected) {
		t.Fatalf("Expected: %s, got: %s", expected, first)
	}
}

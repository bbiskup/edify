package tests

import (
	"bufio"
	sp "edifice/edifact/specparser"
	"reflect"
	"strings"
	"testing"
)

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

	allLines, err := scanner.GetAllSpecLines()
	if err != nil {
		t.Fatalf("Error reading spec lines: %s", err)
	}
	if !reflect.DeepEqual(allLines, expected) {
		t.Fatalf("Expected: %s, got: %s", expected, allLines)
	}
}

func TestSpecScannerFileNotExistent(t *testing.T) {
	_, err := sp.NewSpecScanner("../testdata/specscanner/__NONEXISTANT__")
	if err == nil {
		t.Fatalf("NewSpecScanner should fail with nonexistent file")
	}
}

var specScannerSpecs = []struct {
	inContents string
	expected   [][]string
}{
	{`one

`,
		[][]string{
			{"one"},
		},
	},
	{`one
-------------------------
two`,
		[][]string{
			{"one"},
			{"two"},
		},
	},
	{`one
-------------------------
two

three
four
`,
		[][]string{
			{"one"},
			{"two", "three", "four"},
		},
	},
}

func TestSpecScannerFromReader(t *testing.T) {
	for _, spec := range specScannerSpecs {
		reader := strings.NewReader(spec.inContents)
		bufReader := bufio.NewReader(reader)
		scanner := sp.NewSpecScannerFromReader(bufReader)

		allLines, err := scanner.GetAllSpecLines()
		if err != nil {
			t.Fatalf("Error reading spec lines: %s", err)
		}

		if !reflect.DeepEqual(allLines, spec.expected) {
			t.Fatalf("Expected: %s, got: %s", spec.expected, allLines)
		}
	}
}

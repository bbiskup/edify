package util

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestSpecScannerFromFile(t *testing.T) {
	scanner, err := NewSpecScanner("../../testdata/specscanner/1")
	if err != nil {
		t.Fatalf("Error creating SpecScanner: %s", err)
	}

	expectedHeader := []string{"one"}

	expectedBody := [][]string{
		[]string{"two"},
		[]string{"three"},
	}

	if !reflect.DeepEqual(scanner.HeaderLines, expectedHeader) {
		t.Fatalf("Expected: %s, got: %s", expectedHeader, scanner.HeaderLines)
	}

	allLines, err := scanner.GetAllSpecLines(true)
	if err != nil {
		t.Fatalf("Error reading spec lines: %s", err)
	}

	if !reflect.DeepEqual(allLines, expectedBody) {
		t.Fatalf("Expected: %s, got: %s", expectedBody, allLines)
	}
}

func TestSpecScannerFileNotExistent(t *testing.T) {
	_, err := NewSpecScanner("../../testdata/specscanner/__NONEXISTANT__")
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
		[][]string{},
	},
	{`one
-------------------------
two`,
		[][]string{
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
			{"two", "three", "four"},
		},
	},
	{`one
-------------------------
two
-------------------------
three
four
`,
		[][]string{
			{"two"},
			{"three", "four"},
		},
	},
}

func TestSpecScannerFromReader(t *testing.T) {
	for _, spec := range specScannerSpecs {
		reader := strings.NewReader(spec.inContents)
		bufReader := bufio.NewReader(reader)
		scanner, err := NewSpecScannerFromReader(bufReader)
		if err != nil {
			t.Errorf("Error creating spec scanner from reader: %s", err)
		}

		allLines, err := scanner.GetAllSpecLines(true)
		if err != nil {
			t.Errorf("Error reading spec lines: %s", err)
		}

		if !reflect.DeepEqual(allLines, spec.expected) {
			t.Errorf("Expected: %s, got: %s", spec.expected, allLines)
		}
	}
}

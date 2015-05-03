package specutil

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"
)

func TestSpecScannerFromFile(t *testing.T) {
	scanner, err := NewSpecScanner("../../../testdata/specscanner/1")
	assert.Nil(t, err)
	require.NotNil(t, scanner)

	expectedHeader := []string{"one"}

	expectedBody := [][]string{
		[]string{"two"},
		[]string{"three"},
	}

	assert.True(t, reflect.DeepEqual(scanner.HeaderLines, expectedHeader))

	allLines, err := scanner.GetAllSpecLines(true)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(allLines, expectedBody))
}

func TestSpecScannerFileNotExistent(t *testing.T) {
	scanner, err := NewSpecScanner("../../testdata/specscanner/__NONEXISTANT__")
	assert.NotNil(t, err)
	assert.Nil(t, scanner)
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
		assert.Nil(t, err)

		allLines, err := scanner.GetAllSpecLines(true)
		assert.Nil(t, err)
		assert.True(t, reflect.DeepEqual(allLines, spec.expected))
	}
}

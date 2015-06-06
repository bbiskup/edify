package codes

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFile(t *testing.T) {
	parser := NewCodesSpecParser()
	const fileName = "../../../testdata/UNCL.14B"
	specsMap, err := parser.ParseSpecFile(fileName)
	require.Nil(t, err)

	assert.Equal(t, 272, len(specsMap))
}

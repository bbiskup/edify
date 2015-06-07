package query

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseValidQueryStringOnePartWithIndex(t *testing.T) {
	queryStr := "msg:abc[0]"
	parser, err := NewQueryParser(queryStr)
	assert.Nil(t, err)
	require.NotNil(t, parser)
	assert.Equal(t, 1, len(parser.queryParts))
	assert.Equal(t, queryStr, parser.queryStr)

	part0 := parser.queryParts[0]
	assert.Equal(t, MessageKind, part0.ItemKind)
	assert.Equal(t, "abc", part0.Id)
	assert.Equal(t, 0, part0.Index)
}

func TestParseValidQueryStringOnePartWithoutIndex(t *testing.T) {
	queryStr := "msg:abc"
	parser, err := NewQueryParser(queryStr)
	assert.Nil(t, err)
	require.NotNil(t, parser)
	assert.Equal(t, 1, len(parser.queryParts))
	assert.Equal(t, queryStr, parser.queryStr)

	part0 := parser.queryParts[0]
	assert.Equal(t, MessageKind, part0.ItemKind)
	assert.Equal(t, "abc", part0.Id)
	assert.Equal(t, noIndex, part0.Index)
}

func TestParseValidQueryStringTwoPartsWithoutIndex(t *testing.T) {
	queryStr := "msg:abc|seg:def"
	parser, err := NewQueryParser(queryStr)
	assert.Nil(t, err)
	require.NotNil(t, parser)
	assert.Equal(t, 2, len(parser.queryParts))
	assert.Equal(t, queryStr, parser.queryStr)

	part0 := parser.queryParts[0]
	assert.Equal(t, MessageKind, part0.ItemKind)
	assert.Equal(t, "abc", part0.Id)
	assert.Equal(t, noIndex, part0.Index)

	part1 := parser.queryParts[1]
	assert.Equal(t, SegKind, part1.ItemKind)
	assert.Equal(t, "def", part1.Id)
	assert.Equal(t, noIndex, part1.Index)
}

func TestParseValidQueryStringTwoPartsWithIndex(t *testing.T) {
	queryStr := "msg:abc[2]|seg:def[3]"
	parser, err := NewQueryParser(queryStr)
	assert.Nil(t, err)
	require.NotNil(t, parser)
	assert.Equal(t, 2, len(parser.queryParts))
	assert.Equal(t, queryStr, parser.queryStr)

	part0 := parser.queryParts[0]
	assert.Equal(t, MessageKind, part0.ItemKind)
	assert.Equal(t, "abc", part0.Id)
	assert.Equal(t, 2, part0.Index)

	part1 := parser.queryParts[1]
	assert.Equal(t, SegKind, part1.ItemKind)
	assert.Equal(t, "def", part1.Id)
	assert.Equal(t, 3, part1.Index)
}

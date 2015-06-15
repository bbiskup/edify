package specparser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFullSpecParser(t *testing.T) {
	parser, err := NewFullSpecParser("14B", "../../../testdata/d14b")
	assert.Nil(t, err)
	require.NotNil(t, parser)

	err = parser.Parse()
	assert.Nil(t, err)
	assert.NotNil(t, parser.CodesSpecs)
	assert.NotNil(t, parser.SimpleDataElemSpecs)
	assert.NotNil(t, parser.CompositeDataElemSpecs)

	fmt.Printf("p: %#v", parser)
}

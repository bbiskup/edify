package specparser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFullSpecParser(t *testing.T) {
	// TODO provide permanent test data
	parser, err := NewFullSpecParser("14B", "../../../testdata/d14b")
	assert.Nil(t, err)
	require.NotNil(t, parser)

	err = parser.Parse()
	assert.Nil(t, err)
	fmt.Printf("p: %#v", parser)
}

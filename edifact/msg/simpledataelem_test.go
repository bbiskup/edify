package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElem(t *testing.T) {
	elem := NewSimpleDataElem("ABC", "DEF")
	assert.Equal(t, "ABC", elem.Id())
	assert.Equal(t, "DEF", elem.Value)
	assert.Equal(t, "SimpleDataElem ABC: 'DEF'", elem.String())
}

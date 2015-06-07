package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompositeDataElem(t *testing.T) {
	elem := NewCompositeDataElem("ABC")
	assert.Equal(t, "ABC", elem.Id())
	assert.Equal(t, "CompositeDataElem ABC (0 simple data elems)", elem.String())
}

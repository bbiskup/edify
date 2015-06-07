package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElemString(t *testing.T) {
	elem := NewDataElem([]string{"testValue"})
	assert.Equal(t, "DataElem _no_id 'testValue'", elem.String())
	assert.Equal(t, true, elem.IsSimple())
}

func TestCompositeDataElemString(t *testing.T) {
	elem := NewDataElem([]string{"testValue1", "testvalue2"})
	assert.Equal(t, "DataElem _no_id 'testValue1' 'testvalue2'", elem.String())
	assert.Equal(t, false, elem.IsSimple())
}

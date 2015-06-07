package rawmsg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElemString(t *testing.T) {
	elem := NewRawDataElem([]string{"testValue"})
	assert.Equal(t, "RawDataElem 'testValue'", elem.String())
	assert.Equal(t, true, elem.IsSimple())
}

func TestCompositeDataElemString(t *testing.T) {
	elem := NewRawDataElem([]string{"testValue1", "testvalue2"})
	assert.Equal(t, "RawDataElem 'testValue1' 'testvalue2'", elem.String())
	assert.Equal(t, false, elem.IsSimple())
}

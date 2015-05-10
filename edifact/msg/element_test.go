package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElementString(t *testing.T) {
	elem := NewElement([]string{"testValue"})
	assert.Equal(t, "DataElement 'testValue'", elem.String())
}

func TestCompositeDataElementString(t *testing.T) {
	elem := NewElement([]string{"testValue1", "testvalue2"})
	assert.Equal(t, "DataElement 'testValue1' 'testvalue2'", elem.String())
}

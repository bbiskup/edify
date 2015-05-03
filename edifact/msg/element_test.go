package edifact

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleDataElementString(t *testing.T) {
	elem := NewElement("testName", "testValue")
	assert.Equal(t, "testName testValue", elem.String())
}

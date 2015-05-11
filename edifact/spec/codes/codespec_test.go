package codes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeSpecString(t *testing.T) {
	spec := NewCodeSpec("10", "testname", "testdescr")
	assert.Equal(t, "10 testname testdescr", spec.String())
}

package codes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodesSpec(t *testing.T) {
	spec := NewCodesSpec("10", "codesspec_name", "codespec_descr", []*CodeSpec{
		NewCodeSpec("10", "codespec_name", "codespec_descr"),
	})
	expected := "10 codesspec_name codespe...\n\t10 codespec_name codespe..."
	assert.Equal(t, expected, spec.String())
}

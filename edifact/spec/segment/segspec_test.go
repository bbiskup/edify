package segment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSegSpec(t *testing.T) {
	spec := NewSegSpec("ADR", "ADDRESS", "To specify an address.", nil)

	assert.Equal(t, "Seg ADR/ADDRESS (0 data elems)", spec.String(), "Incorrect String()")
	assert.Equal(t, 0, spec.NumLeadingMandDataElems())
}

func TestParseHeader(t *testing.T) {
	spec := NewSegSpecParser(nil, nil)
	id, name, err := spec.parseHeader("       ADR  ADDRESS")

	assert.Nil(t, err)
	assert.Equal(t, "ADR", id)
	assert.Equal(t, "ADDRESS", name)
}

const funcStr = `       Function: To provide information concerning pricing
                 related to class of trade, price multiplier, and
                 reason for change.`
const expectedFun = `To provide information concerning pricing related to class of trade, price multiplier, and reason for change.`

func TestParseFunction(t *testing.T) {
	spec := NewSegSpecParser(nil, nil)
	funcLines := strings.Split(funcStr, "\n")
	fun, err := spec.parseFunction(funcLines)
	assert.Nil(t, err)
	assert.Equal(t, expectedFun, fun)
}

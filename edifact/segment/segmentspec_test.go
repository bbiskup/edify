package segment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSegmentSpec(t *testing.T) {
	spec := NewSegmentSpec("ADR", "ADDRESS", "To specify an address.", nil)

	assert.Equal(t, "Segment ADR/ADDRESS (0 data elems)", spec.String(), "Incorrect String()")
}

func TestParseHeader(t *testing.T) {
	spec := NewSegmentSpecParser(nil, nil)
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
	spec := NewSegmentSpecParser(nil, nil)
	funcLines := strings.Split(funcStr, "\n")
	fun, err := spec.parseFunction(funcLines)
	assert.Nil(t, err)
	assert.Equal(t, expectedFun, fun)
}

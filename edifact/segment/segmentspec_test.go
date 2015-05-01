package segment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSegmentSpec(t *testing.T) {
	spec := NewSegmentSpec("ADR", "ADDRESS", "To specify an address.", nil)

	const expectedStr = "Segment ADR/ADDRESS (0 data elems)"
	specStr := spec.String()

	if specStr != expectedStr {
		t.Errorf("Expected: %s got: %s", expectedStr, specStr)
	}
}

func TestParseHeader(t *testing.T) {
	spec := NewSegmentSpecParser(nil, nil)

	id, name, err := spec.parseHeader("       ADR  ADDRESS")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	assert.Equal(t, "ADR", id, "Incorrect id")
	assert.Equal(t, "ADDRESS", name, "Incorrect name")
}

const funcStr = `       Function: To provide information concerning pricing
                 related to class of trade, price multiplier, and
                 reason for change.`
const expectedFun = `To provide information concerning pricing related to class of trade, price multiplier, and reason for change.`

func TestParseFunction(t *testing.T) {
	spec := NewSegmentSpecParser(nil, nil)
	funcLines := strings.Split(funcStr, "\n")
	fun, err := spec.parseFunction(funcLines)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	assert.Equal(t, expectedFun, fun, "Incorrect fun")
}

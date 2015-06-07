package segment

import (
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSegSpec(t *testing.T) {
	spec := NewSegSpec("ADR", "ADDRESS", "To specify an address.", nil)

	assert.Equal(t, "Seg ADR/ADDRESS (0 data elems)", spec.String(), "Incorrect String()")
	assert.Equal(t, 0, spec.NumLeadingMandDataElems())
}

func TestParseHeader(t *testing.T) {
	parser := NewSegSpecParser(nil, nil)
	id, name, err := parser.parseHeader("       ADR  ADDRESS")

	assert.Nil(t, err)
	assert.Equal(t, "ADR", id)
	assert.Equal(t, "ADDRESS", name)
}

const funcStr = `       Function: To provide information concerning pricing
                 related to class of trade, price multiplier, and
                 reason for change.`
const expectedFun = `To provide information concerning pricing related to class of trade, price multiplier, and reason for change.`

func TestParseFunction(t *testing.T) {
	parser := NewSegSpecParser(nil, nil)
	funcLines := strings.Split(funcStr, "\n")
	fun, err := parser.parseFunction(funcLines)
	assert.Nil(t, err)
	assert.Equal(t, expectedFun, fun)
}

func TestNumLeadingMandDataElems(t *testing.T) {
	sd1, err := dsp.NewSimpleDataElemSpec(
		"d1", "d1_name", "d1_descr",
		dsp.NewRepr(dsp.AlphaNum, true, 3), nil)
	require.Nil(t, err)

	sd2, err := dsp.NewSimpleDataElemSpec(
		"d1", "d1_name", "d1_descr",
		dsp.NewRepr(dsp.AlphaNum, false, 3), nil)
	require.Nil(t, err)

	spec := NewSegSpec("ABC", "ABC_name", "ABC_descr",
		[]*SegDataElemSpec{NewSegDataElemSpec(sd1, 1, true), NewSegDataElemSpec(sd2, 1, false)})
	assert.Equal(t, 1, spec.NumLeadingMandDataElems())
}

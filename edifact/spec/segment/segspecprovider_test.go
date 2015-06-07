package segment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsUnValidatedSegment(t *testing.T) {
	assert.Equal(t, false, IsUnValidatedSegment("ADR"))
	assert.Equal(t, true, IsUnValidatedSegment("UNH"))
	assert.Equal(t, true, IsUnValidatedSegment("UNT"))
	assert.Equal(t, true, IsUnValidatedSegment("UNS"))
	assert.Equal(t, true, IsUnValidatedSegment("UGH"))
	assert.Equal(t, true, IsUnValidatedSegment("UGT"))
}

func TestSegSpecProvider(t *testing.T) {
	segSpecs := SegSpecMap{
		"ADR": NewSegSpec("ADR", "adr_name", "adr_func", nil),
	}
	prov := NewSegSpecProviderImpl(segSpecs)
	assert.Equal(t, 1, prov.Len())

	adr := prov.Get("ADR")
	require.NotNil(t, adr)
	assert.Equal(t, "ADR", adr.Id)
}

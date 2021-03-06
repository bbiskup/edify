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

func getSegSpecProviderImpl(t *testing.T) *SegSpecProviderImpl {
	segSpecs := SegSpecMap{
		"ADR": NewSegSpec("ADR", "adr_name", "adr_func", nil),
	}
	prov := NewSegSpecProviderImpl(segSpecs)
	return prov
}

func TestSegSpecProviderGetLen(t *testing.T) {
	prov := getSegSpecProviderImpl(t)
	assert.Equal(t, 1, prov.Len())
	adr := prov.Get("ADR")
	require.NotNil(t, adr)
	assert.Equal(t, "ADR", adr.Id)
}

func TestSegSpecProviderGetUnValidated(t *testing.T) {
	prov := getSegSpecProviderImpl(t)
	adr := prov.Get("UNH")
	require.NotNil(t, adr)
	assert.Equal(t, "UNH", adr.Id)
}

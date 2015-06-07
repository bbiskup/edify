package segment

import (
	"github.com/stretchr/testify/assert"
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

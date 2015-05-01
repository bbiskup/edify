package segment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const dataElemSpecStr = "010    C817 ADDRESS USAGE                              C    1"

func TestParseDataElemSpec(t *testing.T) {
	p := NewSegmentSpecParser(nil, nil)

	pos, id, dataElementKind, count, isMandatory, err := p.parseDataElemSpec(dataElemSpecStr)
	assert.Nil(t, err)
	assert.Equal(t, 10, pos)
	assert.Equal(t, "C817", id)
	assert.Equal(t, 1, count)
	assert.False(t, isMandatory)
	assert.Equal(t, Composite, dataElementKind)
}

const segmentSpec = `
       CDI  PHYSICAL OR LOGICAL STATE

       Function: To describe a physical or logical state.

010    7001 PHYSICAL OR LOGICAL STATE TYPE CODE
            QUALIFIER                                  M    1 an..3

020    C564 PHYSICAL OR LOGICAL STATE INFORMATION      M    1
       7007  Physical or logical state description
             code                                      C      an..3
       1131  Code list identification code             C      an..17
       3055  Code list responsible agency code         C      an..3
       7006  Physical or logical state description     C      an..70
`

func TestParseSpec(t *testing.T) {
	p := NewSegmentSpecParser(nil, nil)
	specLines := strings.Split(segmentSpec, "\n")
	segmentSpec, err := p.ParseSpec(specLines)

	assert.Nil(t, err)
	assert.Equal(t, "CDI", segmentSpec.Id)
	assert.Equal(t, "PHYSICAL OR LOGICAL STATE", segmentSpec.Name)
	assert.Equal(t, "To describe a physical or logical state.", segmentSpec.Function)

	lenSegmentDataElementSpecs := len(segmentSpec.SegmentDataElementSpecs)
	assert.Equal(t, 2, lenSegmentDataElementSpecs)
}

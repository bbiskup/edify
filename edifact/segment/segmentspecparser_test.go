package segment

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const dataElemSpecStr = "010    C817 ADDRESS USAGE                              C    1"

func TestParseDataElemSpec(t *testing.T) {
	p := NewSegmentSpecParser(nil, nil)

	pos, id, dataElementKind, count, isMandatory, err := p.parseDataElemSpec(dataElemSpecStr)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error parsing data elem spec: %s", err))
	}

	assert.Equal(t, 10, pos, "Incorrect position")
	assert.Equal(t, "C817", id, "Incorrect ID")
	assert.Equal(t, 1, count, "Incorrect count")
	assert.False(t, isMandatory, "Should be conditional")
	assert.Equal(t, Composite, dataElementKind, "Should be composite")
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

	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if segmentSpec.Id != "CDI" {
		t.Errorf("Id should be CDI; was %s", segmentSpec.Id)
	}

	if segmentSpec.Name != "PHYSICAL OR LOGICAL STATE" {
		t.Errorf("Name should be 'PHYSICAL OR LOGICAL STATE'; was %s", segmentSpec.Name)
	}

	if segmentSpec.Function != "To describe a physical or logical state." {
		t.Errorf("Name should be 'To describe a physical or logical state.'; was %s", segmentSpec.Name)
	}

	lenSegmentDataElementSpecs := len(segmentSpec.SegmentDataElementSpecs)
	if lenSegmentDataElementSpecs != 2 {
		t.Errorf("Expected 2 data element specs; got %d", lenSegmentDataElementSpecs)
	}
}

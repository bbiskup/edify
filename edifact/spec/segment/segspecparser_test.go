package segment

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const dataElemSpecStr = "010    C817 ADDRESS USAGE                              C    1"

func TestParseDataElemSpec(t *testing.T) {
	p := NewSegSpecParser(nil, nil)

	pos, id, dataElemKind, count, isMandatory, err := p.parseDataElemSpec(dataElemSpecStr)
	assert.Nil(t, err)
	assert.Equal(t, 10, pos)
	assert.Equal(t, "C817", id)
	assert.Equal(t, 1, count)
	assert.False(t, isMandatory)
	assert.Equal(t, Composite, dataElemKind)
}

const segSpec = `
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
	p := NewSegSpecParser(nil, nil)
	specLines := strings.Split(segSpec, "\n")
	segSpec, err := p.ParseSpec(specLines)

	assert.Nil(t, err)
	assert.Equal(t, "CDI", segSpec.Id)
	assert.Equal(t, "PHYSICAL OR LOGICAL STATE", segSpec.Name)
	assert.Equal(t, "To describe a physical or logical state.", segSpec.Function)

	lenSegDataElemSpecs := len(segSpec.SegDataElemSpecs)
	assert.Equal(t, 2, lenSegDataElemSpecs)
}

const segSpecWithNote = `
       ARR  ARRAY INFORMATION

       Function: To contain the data in an array.

010    C778 POSITION IDENTIFICATION                    C    1
       7164  Hierarchical structure level identifier   C      an..35
       1050  Sequence position identifier              C      an..10

020    C770 ARRAY CELL DETAILS                         C    1
       9424  Array cell data description               C      an..512

       Note: 
            The composite C770 - array cell details - occurs
            9,999 times in the segment. The use of the ARR
            segment is restricted to be used only with Version 3
            of ISO-9735.
            The component 9424 - array cell information - occurs
            100 times in the composite C770. The use of C770 is
            restricted to be used only with the ARR segment
            within Version 3 of ISO-9735.
`

func TestParseSpecWithNode(t *testing.T) {
	p := NewSegSpecParser(nil, nil)
	specLines := strings.Split(segSpecWithNote, "\n")
	segSpec, err := p.ParseSpec(specLines)

	assert.Nil(t, err)
	assert.Equal(t, "ARR", segSpec.Id)
	lenSegDataElemSpecs := len(segSpec.SegDataElemSpecs)
	assert.Equal(t, 2, lenSegDataElemSpecs)
}

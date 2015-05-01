package dataelement

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func TestSpecParser(t *testing.T) {
	// TODO will fail; does not contain required keys
	p := NewSimpleDataElementSpecParser(fixtureTextCodesSpecMap())
	_, err := p.ParseSpecFile("../../testdata/EDED.14B_short")
	assert.Nil(t, err)
}

const specLines = `
     1000  Document name                                           [B]

     Desc: Name of a document.

     Repr: an..35
`

func TestParseSpecLines(t *testing.T) {
	p := NewSimpleDataElementSpecParser(fixtureTextCodesSpecMap())
	res, err := p.ParseSpec(strings.Split(specLines, "\n"))
	assert.Nil(t, err)
	log.Printf("res: %s", res)
}

func BenchmarkParseSpecLines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewSimpleDataElementSpecParser(fixtureTextCodesSpecMap())
		specs, err := p.ParseSpecFile("../../testdata/EDED.14B")
		assert.Nil(b, err)
		log.Printf("Parsed %d specs\n", len(specs))
	}
}

package validation

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// Check feasibility of using regexp matching to find valid sequences
// segments according to a message spec
func createLargeRegexpStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	for i := 0; i < 1000; i++ {

		// CCC with count 1 makes huge difference: 1000 x faster
		// than pattern with only optional AAA and BBB
		buffer.WriteString("(AAA:){0,3}(BBB:){0,2}(CCC:)")
	}
	buffer.WriteString("$")
	return buffer.String()
}

func BenchmarkCompileLargeRegexp(b *testing.B) {
	patStr := createLargeRegexpStr()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := regexp.MustCompile(patStr)

		if i == 0 {
			assert.NotNil(b, r)
		}
	}
}

func BenchmarkMatchLargeRegexp(b *testing.B) {
	patStr := createLargeRegexpStr()
	r := regexp.MustCompile(patStr)
	assert.NotNil(b, r)

	numSegments := 1000

	var buffer bytes.Buffer
	for i := 0; i < numSegments; i++ {
		var next string
		if i%2 == 0 {
			next = "AAA:AAA:AAA:"
		} else {
			next = "BBB:BBB:"
		}
		buffer.WriteString(next)
		buffer.WriteString("CCC:")
	}
	strToMatch := buffer.String()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := r.FindStringSubmatch(strToMatch)

		if i == 0 {
			assert.NotNil(b, m)

			// + 1 for group 0 (full match)
			assert.Equal(b, numSegments*3+1, len(m))
		}
	}
}

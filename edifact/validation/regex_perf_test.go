package validation

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func BenchmarkCompileLargeRegexp(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < 1000; i++ {
		buffer.WriteString("(AAA:{1,3})(BBB:{0,2})")
	}
	patStr := buffer.String()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := regexp.MustCompile(patStr)
		assert.NotNil(b, r)
	}
}

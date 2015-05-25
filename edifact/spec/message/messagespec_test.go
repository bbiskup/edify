package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var partsSpec = []struct {
	parts    []MessageSpecPart
	expected string
}{
	{
		parts:    []MessageSpecPart{},
		expected: "Message testid (testname testrelease): 0 parts",
	},
}

func TestString(t *testing.T) {
	for _, spec := range partsSpec {
		s := NewMessageSpec(
			"testid", "testname",
			"testversion", "testrelease",
			"testcontragency", "testrevision",
			time.Date(2015, time.January, 15, 0, 0, 0, 0, time.UTC),
			"testsource",
			[]MessageSpecPart{},
		)

		assert.Equal(t, spec.expected, s.String())
	}
}

package message

import (
	"github.com/bbiskup/edify/edifact/spec/segment"
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

	{
		parts: []MessageSpecPart{
			NewMessageSpecSegmentPart(
				segment.NewSegmentSpec("UNH", "testname1", "testfunc1", nil), 1, false, nil),
			NewMessageSpecSegmentPart(
				segment.NewSegmentSpec("BGM", "testname1", "testfunc1", nil), 1, false, nil),
			NewMessageSpecSegmentPart(
				segment.NewSegmentSpec("UNT", "testname1", "testfunc1", nil), 1, false, nil),
		},
		expected: "Message testid (testname testrelease): 3 parts - UNH, BGM, UNT",
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
			spec.parts,
		)

		assert.Equal(t, spec.expected, s.String())
	}
}

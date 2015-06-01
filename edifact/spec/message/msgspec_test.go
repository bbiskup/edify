package message

import (
	"github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var partsSpec = []struct {
	parts    []MsgSpecPart
	expected string
}{
	{
		parts:    []MsgSpecPart{},
		expected: "Message testid (testname testrelease): 0 parts",
	},

	{
		parts: []MsgSpecPart{
			NewMsgSpecSegPart(
				segment.NewSegSpec("UNH", "testname1", "testfunc1", nil), 1, false, nil),
			NewMsgSpecSegPart(
				segment.NewSegSpec("BGM", "testname1", "testfunc1", nil), 1, false, nil),
			NewMsgSpecSegPart(
				segment.NewSegSpec("UNT", "testname1", "testfunc1", nil), 1, false, nil),
		},
		expected: "Message testid (testname testrelease): 3 parts - UNH, BGM, UNT",
	},
}

func TestString(t *testing.T) {
	for _, spec := range partsSpec {
		s := NewMsgSpec(
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

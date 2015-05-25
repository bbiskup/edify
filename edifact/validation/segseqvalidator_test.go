package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var authorSegSeqSpec = []struct {
	descr       string
	segmentIDs  []string
	expectError bool
	errorKind   SegSeqErrorKind
}{
	{
		"No segments at all",
		[]string{}, true, noSegments,
	},

	{
		"Missing mandatory segments",
		[]string{
			"UNH", // no BGM
		}, true, missingMandatorySegment,
	},

	{
		"Max. repeat count of mandatory segment exeeded",
		[]string{
			"UNH", "UNH", // max. repeat count is 1
		}, true, maxSegmentRepeatCountExceeded},

	{
		"Max. repeat count of optional segment exeeded",
		[]string{
			"UNH", "BGM", "DTM" /* max. repeat count is 1 */, "DTM",
		}, true, maxSegmentRepeatCountExceeded},

	{"Optional segment in incorrect position",
		[]string{
			"UNH",
			"DTM" /* Should appear after BGM */, "BGM", "UNT",
		}, true, unexpectedSegment,
	},

	{"Optional segment in incorrect position",
		[]string{
			"DTM", "UNH", "BGM", "UNT",
		}, true, missingMandatorySegment,
	},

	{"Missing mandatory group 4",
		[]string{
			"UNH", "BGM", "DTM" /* optional */, "UNT",
		}, true, missingGroup,
	},

	{"minimal message (only mandatory segments)",
		[]string{
			"UNH", "BGM" /* Group 4 */, "LIN",

			"UNT",
		}, false, "",
	},

	{
		"Mostly mandatory",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS", // both conditional
			// Group 4
			"LIN",
			"UNT",
		}, false, ""},

	{
		"Mostly mandatory; one conditional group",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 1
			"LIN",
			// Group 2
			"FII", "CTA", "COM",

			"UNT",
		}, false, ""},

	{
		"Some repeat counts > 1",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 4
			"LIN", "LIN", "LIN", "LIN",
			// Group 7
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",

			"UNT",
		}, false, ""},

	/*

		{
			"Some repeat counts > 1",
			[]string{
				"UNH", "BGM",
				"DTM", "BUS",
				// Group 4
				"LIN", "LIN", "LIN", "LIN",
				// Group 7
				"FII", "CTA", "COM", "COM", "COM",
				"FII", "CTA", "COM", "COM", "COM",

				"UNT",
			}, false},


		{
			"group 7 repeated too often",
			[]string{
				"UNH", "BGM",
				"DTM", "BUS",
				// Group 4
				"LIN", "LIN", "LIN", "LIN",
				// Group 7
				"FII", "CTA", "COM", "COM", "COM",
				"FII", "CTA", "COM", "COM", "COM",
				"FII", "CTA", "COM", "COM", "COM",

				"UNT",
			}, true},*/
}

func TestSegSeqValidator1(t *testing.T) {
	msgSpec := getMessageSpec("AUTHOR_D.14B")

	for _, spec := range authorSegSeqSpec {
		fmt.Printf(">>>>>>>>>>>>>>>>>>> spec: %#v\n", spec)
		validator, err := NewSegSeqValidator(msgSpec)
		require.Nil(t, err)
		require.NotNil(t, validator)
		segments := mapToSegments(spec.segmentIDs)
		require.NotNil(t, segments)
		message := msg.NewMessage("AUTHOR", segments)
		fmt.Printf("Validating message: %s", message)
		err = validator.Validate(message)

		if spec.expectError {
			assert.NotNil(t, err)
			fmt.Printf("Expected error was: %s\n", err)
			err, ok := err.(SegSeqError)
			assert.True(t, ok)
			require.Equal(t, spec.errorKind, err.kind)
		} else {
			require.Nil(t, err)
		}
	}
}

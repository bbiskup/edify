package validation

import (
	"github.com/bbiskup/edify/edifact/spec/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getMessageSpec() *message.MessageSpec {
	parser := message.NewMessageSpecParser(&message.MockSegmentSpecProviderImpl{})
	messageSpec, err := parser.ParseSpecFile("../../testdata/AUTHOR_D.14B")
	if err != nil {
		panic("spec is nil")
	}
	return messageSpec
}

func TestMessageValidator(t *testing.T) {
	validator, err := NewMessageValidator(getMessageSpec())
	require.Nil(t, err)
	require.NotNil(t, validator)

	expected := "^(UNH:)(BGM:)(DTM:){0,1}(BUS:){0,1}((RFF:)(DTM:){0,1}){0,2}((FII:)(CTA:){0,1}(COM:){0,5}){0,5}((NAD:)(CTA:){0,1}(COM:){0,5}){0,3}((LIN:)((RFF:)(DTM:){0,1}){0,5}((SEQ:)(GEI:)(DTM:){0,2}(MOA:){0,1}(DOC:){0,5}){0,}((FII:)(CTA:){0,1}(COM:){0,5}){0,2}((NAD:)(CTA:){0,1}(COM:){0,5}){0,2}){1,}(CNT:){0,5}((AUT:)(DTM:){0,1}){0,5}(UNT:)$"
	assert.Equal(t, expected, validator.segmentValidationRegexpStr)
}

var segListStrSpec = []struct {
	segmentIDs []string
	expected   string
}{
	{[]string{}, ""},
	{[]string{"AAA"}, "AAA:"},
	{[]string{"AAA", "BBB"}, "AAA:BBB:"},
}

func TestBuildSegmentListStr(t *testing.T) {
	for _, spec := range segListStrSpec {
		result := buildSegmentListStr(spec.segmentIDs)
		assert.Equal(t, spec.expected, result)
	}
}

var authorSegSeqSpec = []struct {
	segmentIDs  []string
	valid       bool
	expectError bool
}{
	// minimal message (only mandatory segments)
	{[]string{
		"UNH", "BGM",
		// Group 1
		"LIN",
		"UNT",
	}, true, false},

	// Mostly mandatory
	{[]string{
		"UNH", "BGM",
		"DTM", "BUS", // both conditional
		// Group 4
		"LIN",
		"UNT",
	}, true, false},

	// Mostly mandatory; one conditional group
	{[]string{
		"UNH", "BGM",
		"DTM", "BUS",
		// Group 1
		"LIN",
		// Group 2
		"FII", "CTA", "COM",

		"UNT",
	}, true, false},

	// Some repeat counts > 1
	{[]string{
		"UNH", "BGM",
		"DTM", "BUS",
		// Group 4
		"LIN", "LIN", "LIN", "LIN",
		// Group 7
		"FII", "CTA", "COM", "COM", "COM",
		"FII", "CTA", "COM", "COM", "COM",

		"UNT",
	}, true, false},

	// No segments at all
	{[]string{}, false, true},

	// Missing mandatory segments
	{[]string{"UNH"}, false, false},

	// First mandatory segment repeated too often
	{[]string{
		"UNH", "UNH", "BGM",
		// Group 1
		"LIN",
		"UNT",
	}, false, false},

	// group 7 repeated too often
	{[]string{
		"UNH", "BGM",
		"DTM", "BUS",
		// Group 4
		"LIN", "LIN", "LIN", "LIN",
		// Group 7
		"FII", "CTA", "COM", "COM", "COM",
		"FII", "CTA", "COM", "COM", "COM",
		"FII", "CTA", "COM", "COM", "COM",

		"UNT",
	}, false, false},
}

func TestValidateSegmentList(t *testing.T) {
	validator, err := NewMessageValidator(getMessageSpec())
	require.Nil(t, err)
	// fmt.Printf("regexp str %s", validator.segmentValidationRegexpStr)
	for _, spec := range authorSegSeqSpec {
		// fmt.Printf("spec: %#v\n", spec)
		result, err := validator.ValidateSegmentList(spec.segmentIDs)
		assert.Equal(t, spec.valid, result)
		assert.Equal(t, spec.expectError, err != nil)
	}
}

// Benchmark creation of validation regexp
func BenchmarkNewMessageValidator(b *testing.B) {
	messageSpec := getMessageSpec()
	for i := 0; i < b.N; i++ {
		validator, err := NewMessageValidator(messageSpec)
		require.Nil(b, err)
		require.NotNil(b, validator)
	}
}

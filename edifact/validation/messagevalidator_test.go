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

// Benchmark creation of validation regexp
func BenchmarkNewMessageValidator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validator, err := NewMessageValidator(getMessageSpec())
		require.Nil(b, err)
		require.NotNil(b, validator)
	}
}

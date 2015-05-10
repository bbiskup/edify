package validation

import (
	"github.com/bbiskup/edify/edifact/spec/message"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageValidator(t *testing.T) {
	parser := message.NewMessageSpecParser(&message.MockSegmentSpecProviderImpl{})
	messageSpec, err := parser.ParseSpecFile("../../testdata/AUTHOR_D.14B")
	require.Nil(t, err)
	require.NotNil(t, messageSpec)
	validator, err := NewMessageValidator(messageSpec)
	require.Nil(t, err)
	require.NotNil(t, validator)
}

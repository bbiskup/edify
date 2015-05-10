package validation

import (
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
)

type MessageValidator struct {
	messageSpec msgspec.MessageSpec
}

func (v *MessageValidator) String() string {
	return fmt.Sprintf("MessageValidator (%s)", v.String())
}

func (v *MessageValidator) Validate(message msg.Message) error {
	panic("Not implemented")
}

func NewMessageValidator(messageSpec msgspec.MessageSpec) *MessageValidator {
	return &MessageValidator{messageSpec}
}

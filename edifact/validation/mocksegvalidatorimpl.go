package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
)

// For testing purposes.
// - Considers every raw seg as valid (does not care about data elements)
// - Creates a new Seg without data elements
type MockSegValidatorImpl struct {
}

func (v *MockSegValidatorImpl) Validate(rawSeg *rawmsg.RawSeg) (*msg.Seg, error) {
	return msg.NewSeg(rawSeg.Id()), nil
}

package validation

import (
	"errors"
	"github.com/bbiskup/edify/edifact/msg"
	//msp "github.com/bbiskup/edify/edifact/spec/message"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
)

// Determines message type from UNH segment.
func getMsgTypeFromUNH(seg *msg.Seg) (msgName string, err error) {
	// Since knowing the message
	// type is a prerequisite to validating a raw message and constructing a
	// nested message from it, this method cannot use the  query mechanism
	if len(seg.Elems) < 2 {
		return "", errors.New("Too few data elements; no message type")
	}
	msgTypeElem := seg.Elems[1]
	if len(msgTypeElem.Values) < 1 {
		return "", errors.New(
			"Too few component elements in msg type composite element; no message type")
	}
	return msgTypeElem.Values[0], nil
}

// Validates an entire message:
// - correctness of segment sequence
// - correctness/completeness of data elements
type MsgValidator struct {
	segSpecs ssp.SegSpecMap
}

func (v *MsgValidator) Validate(rawMsg *msg.RawMsg) (nestedMsg *msg.NestedMsg, err error) {
	if len(rawMsg.Segs) == 0 {
		return nil, errors.New("No segments")
	}
	unh := rawMsg.Segs[0]
	if unh.Id() != "UNH" {
		return nil, errors.New("Could not find UNH segment")
	}

	panic("Not implemented")
}

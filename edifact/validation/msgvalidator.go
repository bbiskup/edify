package validation

import (
	"errors"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
	//msp "github.com/bbiskup/edify/edifact/spec/message"
	"fmt"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
	"strconv"
)

// Determines message type from UNH segment.
func getMsgTypeFromUNH(rawSeg *rawmsg.RawSeg) (msgName string, err error) {
	// Since knowing the message
	// type is a prerequisite to validating a raw message and constructing a
	// nested message from it, this method cannot use the  query mechanism
	if len(rawSeg.Elems) < 2 {
		return "", errors.New("Too few data elements; no message type")
	}
	msgTypeElem := rawSeg.Elems[1]
	if len(msgTypeElem.Values) < 1 {
		return "", errors.New(
			"Too few component elements in msg type composite element; no message type")
	}
	return msgTypeElem.Values[0], nil
}

// Determines number of segments in message from UNT segment
func getSegCountFromUNT(rawSeg *rawmsg.RawSeg) (segCount int, err error) {
	if len(rawSeg.Elems) < 1 {
		return -1, errors.New("Too few data elements; segment count")
	}
	numSegsElem := rawSeg.Elems[0]
	if len(numSegsElem.Values) < 1 {
		return -1, errors.New(
			"Too few component elements in msg type composite element; no segment count")
	}
	segCountStr := numSegsElem.Values[0]
	return strconv.Atoi(segCountStr)
}

// Validates an entire message:
// - correctness of segment sequence
// - correctness/completeness of data elements
type MsgValidator struct {
	msgSpecs     msp.MsgSpecMap
	segSpecs     ssp.SegSpecProvider
	segValidator SegValidator
}

func (v *MsgValidator) Validate(rawMsg *rawmsg.RawMsg) (nestedMsg *msg.NestedMsg, err error) {
	if len(rawMsg.RawSegs) == 0 {
		return nil, errors.New("No segments")
	}
	unh := rawMsg.RawSegs[0]
	if unh.Id() != "UNH" {
		return nil, errors.New("Could not find UNH segment")
	}
	var msgType string
	if msgType, err = getMsgTypeFromUNH(unh); err != nil {
		return nil, err
	}

	unt := rawMsg.RawSegs[len(rawMsg.RawSegs)-1]
	if unt.Id() != "UNT" {
		return nil, errors.New("Could not find UNT segment")
	}

	var segCount int
	if segCount, err = getSegCountFromUNT(unt); err != nil {
		return nil, err
	}

	log.Printf("Validating message %s (%d segments)", msgType, segCount)

	// Validate segments
	for _, rawSeg := range rawMsg.RawSegs {
		err := v.segValidator.Validate(rawSeg)
		if err != nil {
			return nil, err
		}
	}

	// Validate segment sequence
	msgSpec, ok := v.msgSpecs[msgType]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("No message spec found for message type %s", msgType))
	}
	segSeqValidator := NewSegSeqValidator(msgSpec)
	return segSeqValidator.Validate(rawMsg)
}

func NewMsgValidator(msgSpecs msp.MsgSpecMap, segSpecProvider ssp.SegSpecProvider) *MsgValidator {
	return &MsgValidator{msgSpecs, segSpecProvider, NewSegValidatorImpl(segSpecProvider)}
}

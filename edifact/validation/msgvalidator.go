package validation

import (
	"errors"
	"github.com/bbiskup/edify/edifact/msg"
	//msp "github.com/bbiskup/edify/edifact/spec/message"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
	"strconv"
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

// Determines number of segments in message from UNT segment
func getSegCountFromUNT(seg *msg.Seg) (segCount int, err error) {
	if len(seg.Elems) < 1 {
		return -1, errors.New("Too few data elements; segment count")
	}
	numSegsElem := seg.Elems[0]
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
	var msgType string
	if msgType, err = getMsgTypeFromUNH(unh); err != nil {
		return nil, err
	}

	unt := rawMsg.Segs[len(rawMsg.Segs)-1]
	if unt.Id() != "UNT" {
		return nil, errors.New("Could not find UNT segment")
	}

	var segCount int
	if segCount, err = getSegCountFromUNT(unt); err != nil {
		return nil, err
	}

	log.Printf("Validating message %s (%d segments)", msgType, segCount)

	panic("Not implemented")
}

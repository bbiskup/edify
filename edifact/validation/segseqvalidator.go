package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	"log"
	"strconv"
)

// - Validates a segment sequence according to edmd specs
// - Constructs a valid NestedMsg
//
// The validator is not thread-safe
type SegSeqValidator struct {
	// List of segments to be consumed during validation
	segs    []*msg.Seg
	rawMsg  *msg.RawMsg
	msgSpec *msp.MsgSpec
}

func (v *SegSeqValidator) String() string {
	var segStr string
	if v.segs != nil {
		segStr = strconv.FormatInt(int64(len(v.segs)), 10)
	} else {
		segStr = "-"
	}
	return fmt.Sprintf("SegSeqValidator (msg: %s, segments left: %s)",
		v.msgSpec.Id, segStr)
}

// Remove the current segment from the list of segments under
// validation
func (v *SegSeqValidator) consume() {
	if v.segs == nil || len(v.segs) == 0 {
		panic("consume() called on missing/empty segment list")
	}

	// Remove leading segment and all subsequent segments of same ID
	firstSegID := v.segs[0].Id()
	var cutIndex int
	for index, seg := range v.segs {
		cutIndex = index + 1
		if seg.Id() != firstSegID {
			break
		}
	}
	v.segs = v.segs[cutIndex:]
}

// Returns the next segment of the message under validation
// (if one exists) or panics
func (v *SegSeqValidator) peek() []*msg.Seg {
	if v.segsExhausted() {
		panic("No more segments")
	} else {
		firstSegID := v.segs[0].Id()
		result := []*msg.Seg{}
		for _, seg := range v.segs {
			if seg.Id() == firstSegID {
				result = append(result, seg)
			}
		}
		return result
	}
}

// Returns true if all segments of the message under validation
// have been consumed
func (v *SegSeqValidator) segsExhausted() bool {
	return len(v.segs) == 0
}

// Validates segment groups recursively
// while building nested message
func (v *SegSeqValidator) validateGroup(
	//context *SegSeqGroupContext,
	curMsgSpecSegGrpPart *msp.MsgSpecSegGrpPart,
	curRepSegGrp *msg.RepSegGrp,
) (nestedMsg *msg.NestedMsg, err error) {

	log.Printf("Entering group spec %s", curMsgSpecSegGrpPart.Id())

	for _, specPart := range curMsgSpecSegGrpPart.Children() {
		if v.segsExhausted() {
			return nil, NewSegSeqError(
				missingMandatorySeg, "Segments exhausted")
		}
		segs := v.peek()
		repeatCount := len(segs)
		segID := segs[0].Id()
		log.Printf("Peek: %s (%dx)", segID, repeatCount)

		// Generic error msg
		segErrStr := fmt.Sprintf("Segment %s in group %s",
			segID, curMsgSpecSegGrpPart)

		switch specPart := specPart.(type) {
		case *msp.MsgSpecSegPart:
			if specPart.Id() != segID {
				if specPart.IsMandatory() {
					return nil, NewSegSeqError(missingMandatorySeg, segErrStr)
				}
			} else {
				if repeatCount > specPart.MaxCount() {
					return nil, NewSegSeqError(maxSegRepeatCountExceeded, segErrStr)
				}
				log.Printf("Consuming segment %s", segID)
				v.consume()
				continue
			}
		case *msp.MsgSpecSegGrpPart:
		default:
			panic(fmt.Sprintf("Unsupported type %T", specPart))
		}
	}

	panic("Not implemented")
}

func (v *SegSeqValidator) Validate(rawMsg *msg.RawMsg) (nestedMsg *msg.NestedMsg, err error) {
	if len(rawMsg.Segs) == 0 {
		return nil, NewSegSeqError(noSegs, "")
	}

	v.segs = make([]*msg.Seg, len(rawMsg.Segs))
	// Make a copy so the original msg does not get modified
	copy(v.segs, rawMsg.Segs)

	nestedMsg = msg.NewNestedMsg(v.msgSpec.Name)
	//topLevelContext := NewSegSeqGroupContext(v.msgSpec.TopLevelGroup, nestedMsg.TopLevelRepGrp)

	return v.validateGroup(v.msgSpec.TopLevelGroup, nestedMsg.TopLevelRepGrp)
}

func NewSegSeqValidator(msgSpec *msp.MsgSpec) *SegSeqValidator {
	return &SegSeqValidator{
		segs:    nil,
		rawMsg:  nil,
		msgSpec: msgSpec,
	}
}

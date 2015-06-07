package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	//"log"
	"strconv"
)

// - Validates a segment sequence according to edmd specs
// - Constructs a valid NestedMsg
//
// The validator is not thread-safe
type SegSeqValidator struct {
	// List of segments to be consumed during validation
	rawSegs []*rawmsg.RawSeg
	rawMsg  *rawmsg.RawMsg
	msgSpec *msp.MsgSpec
}

func (v *SegSeqValidator) String() string {
	var segStr string
	if v.rawSegs != nil {
		segStr = strconv.FormatInt(int64(len(v.rawSegs)), 10)
	} else {
		segStr = "-"
	}
	return fmt.Sprintf("SegSeqValidator (msg: %s, segments left: %s)",
		v.msgSpec.Id, segStr)
}

// Remove a single segment from the list of segments
func (v *SegSeqValidator) consumeSingle() {
	// log.Printf("consumeSingle()")
	if v.rawSegs == nil || len(v.rawSegs) == 0 {
		panic("consumeSingle() called on missing/empty segment list")
	}
	v.rawSegs = v.rawSegs[1:]
}

// Remove the current segment from the list of segments under
// validation
func (v *SegSeqValidator) consumeMulti() {
	// log.Printf("consumeMulti()")
	if v.rawSegs == nil || len(v.rawSegs) == 0 {
		panic("consumeMulti() called on missing/empty segment list")
	}

	// Remove leading segment and all subsequent segments of same ID
	firstSegID := v.rawSegs[0].Id()
	var cutIndex int
	for index, seg := range v.rawSegs {
		if seg.Id() != firstSegID {
			break
		}
		cutIndex = index + 1
	}
	v.rawSegs = v.rawSegs[cutIndex:]
}

// Returns the next segment of the message under validation
// (if one exists) or panics
func (v *SegSeqValidator) peek() []*rawmsg.RawSeg {
	if v.segsExhausted() {
		panic("No more segments")
	} else {
		firstSegID := v.rawSegs[0].Id()
		result := []*rawmsg.RawSeg{}
		for _, seg := range v.rawSegs {
			if seg.Id() == firstSegID {
				result = append(result, seg)
			} else {
				break
			}
		}
		return result
	}
}

// Returns true if all segments of the message under validation
// have been consumed
func (v *SegSeqValidator) segsExhausted() bool {
	return len(v.rawSegs) == 0
}

func (v *SegSeqValidator) hasRemainingMandatorySpecs(specIndex int, groupChildren []msp.MsgSpecPart) bool {
	numChildren := len(groupChildren)
	for i := specIndex + 1; i < numChildren; i++ {
		if groupChildren[i].IsMandatory() {
			return true
		}
	}
	return false
}

// Validates segment groups recursively
// while building nested message
func (v *SegSeqValidator) validateGroup(
	//context *SegSeqGroupContext,
	curMsgSpecSegGrpPart *msp.MsgSpecSegGrpPart,
	curRepSegGrp *msg.RepSegGrp,
) error {

	//log.Printf("Entering group spec %s", curMsgSpecSegGrpPart.Name())
	// log.Printf("@BUILD: curRepSegGrp: %s", curRepSegGrp.Dump(0))

	groupRepeatCount := 0
	groupTriggerSegmentID := curMsgSpecSegGrpPart.TriggerSegPart().Id()

GROUPREPEAT:
	for {
		//log.Printf("Repeating group %s # %d",
		//	curMsgSpecSegGrpPart.Name(), groupRepeatCount)

		var segGrp *msg.SegGrp
		if !curRepSegGrp.IsTopLevel() {
			segGrp = msg.NewSegGrp(curMsgSpecSegGrpPart.Name())
			// log.Printf("@BUILD: appending %s", segGrp)
			curRepSegGrp.Append(segGrp)
		} else {
			// log.Printf("@BUILD: not appending to top-level group")
			segGrp = curRepSegGrp.GetSegGrp(0)
		}

		groupRepeatCount++

		for specIndex, specPart := range curMsgSpecSegGrpPart.Children() {
			if v.segsExhausted() {
				return NewSegSeqError(
					missingMandatorySeg, "Segments exhausted")
			}
			rawSegs := v.peek()
			repeatCount := len(rawSegs)
			segID := rawSegs[0].Id()
			// log.Printf("Spec: %s; peek: %s (%dx)",
			//	specPart, segID, repeatCount)

			// Generic error msg
			segErrStr := fmt.Sprintf("%s in %s",
				specPart, curMsgSpecSegGrpPart)

			switch specPart := specPart.(type) {
			case *msp.MsgSpecSegPart:
				if specPart.Id() != segID {
					// log.Printf("unequal spec: %s vs seg: %s", specPart.Id(), segID)
					if specPart.IsMandatory() {
						return NewSegSeqError(missingMandatorySeg, segErrStr)
					} else {
						continue
					}
				}

				// Segments are equal
				if repeatCount > specPart.MaxCount() {
					if !v.hasRemainingMandatorySpecs(specIndex, curMsgSpecSegGrpPart.Children()) {
						// A segment repetition may occur if a group contains
						// of a single mandatory segment, and all optional
						// segments are skipped; e.g. AUTHOR message, group_4: LIN segment
						// log.Printf("repeat count exceeded? repeating group")
						v.consumeSingle()

						newRepSeg := msg.NewRepSeg(rawSegs[0])

						// TODO: validate newly added segment

						// log.Printf("@BUILD: Appending %s to %s", newRepSeg, segGrp)
						segGrp.AppendRepSeg(newRepSeg)

						continue GROUPREPEAT
					} else {
						// log.Printf("There are remaining mandatory segments")
						return NewSegSeqError(maxSegRepeatCountExceeded, segErrStr)
					}
				}

				// Segment matches
				newRepSeg := msg.NewRepSeg(rawSegs...)
				// log.Printf("@BUILD: Appending %s to %s", newRepSeg, segGrp)
				segGrp.AppendRepSeg(newRepSeg)
				v.consumeMulti()
				continue
			case *msp.MsgSpecSegGrpPart:
				triggerSegmentID := specPart.TriggerSegPart().Id()
				if triggerSegmentID != segID {
					// log.Printf("Trigger for group %s not present", specPart.Name())
					if specPart.IsMandatory() {
						return NewSegSeqError(
							missingMandatorySeg,
							fmt.Sprintf("Trigger segment %s for group %s",
								triggerSegmentID, specPart.Name()))
					}
				} else {
					newRepSegGrp := msg.NewRepSegGrp(specPart.Name())
					// log.Printf("@BUILD appending %s to %s", newRepSegGrp, segGrp)
					segGrp.AppendRepSegGrp(newRepSegGrp)
					if err := v.validateGroup(specPart, newRepSegGrp); err != nil {
						return err
					}
				}

			default:
				panic(fmt.Sprintf("Unsupported type %T", specPart))
			}
		}
		if v.segsExhausted() {
			return nil
		}
		if groupRepeatCount > curMsgSpecSegGrpPart.MaxCount() {
			return NewSegSeqError(
				maxGroupRepeatCountExceeded,
				fmt.Sprintf("Group %s", curMsgSpecSegGrpPart.Name()))
		}

		if v.peek()[0].Id() != groupTriggerSegmentID {
			break
		}
	}
	//log.Printf("Leaving group spec %s", curMsgSpecSegGrpPart.Name())
	return nil
}

func (v *SegSeqValidator) Validate(rawMsg *rawmsg.RawMsg) (nestedMsg *msg.NestedMsg, err error) {
	if len(rawMsg.RawSegs) == 0 {
		return nil, NewSegSeqError(noSegs, "")
	}

	v.rawSegs = make([]*rawmsg.RawSeg, len(rawMsg.RawSegs))
	// Make a copy so the original msg does not get modified
	copy(v.rawSegs, rawMsg.RawSegs)

	nestedMsg = msg.NewNestedMsg(v.msgSpec.Name)
	//topLevelContext := NewSegSeqGroupContext(v.msgSpec.TopLevelGroup, nestedMsg.TopLevelRepGrp)

	if err := v.validateGroup(v.msgSpec.TopLevelGroup, nestedMsg.TopLevelRepGrp); err != nil {
		return nil, err
	} else {
		return nestedMsg, nil
	}
}

func NewSegSeqValidator(msgSpec *msp.MsgSpec) *SegSeqValidator {
	return &SegSeqValidator{
		rawSegs: nil,
		rawMsg:  nil,
		msgSpec: msgSpec,
	}
}

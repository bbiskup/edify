package validation

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	msp "github.com/bbiskup/edify/edifact/spec/message"
	"github.com/bbiskup/edify/edifact/util"
	"log"
)

// Validates segment sequence
// builds structure for navigation/query
type SegSeqValidator struct {
	// index in raw message
	currentSegIndex int

	// Specification for message type under validation
	msgSpec *msp.MsgSpec

	// ID of the segment specification currently checked
	currentSegSpecID string
	state            SegSeqState

	// Id of current raw message segment
	currentSegID string

	// Manages nesting of groups
	groupStack *util.Stack

	nestedMsgBuilder *NestedMsgBuilder
}

// There is always a current group context (the validator starts
// out with a top-level group context)
func (s *SegSeqValidator) currentGroupContext() *SegSeqGroupContext {
	result := (s.groupStack).Peek().(*SegSeqGroupContext)
	return result
}

func (s *SegSeqValidator) createError(kind SegSeqErrKind, msg string) error {
	return NewSegSeqError(
		kind, fmt.Sprintf("Error at segment #%d (%s)",
			s.currentSegIndex, msg))
}

// Sets new state and logs state transition
func (s *SegSeqValidator) setNewState(newState SegSeqState) {
	log.Printf("State transition %s --> %s", s.state, newState)
	s.state = newState
}

func (s *SegSeqValidator) handleRepSeg(segment *msg.Seg) error {
	log.Printf("handleRepSeg %s", segment)
	gc := s.currentGroupContext()
	gc.segmentRepeatCount++
	maxCount := s.getCurrentMsgSpecPart().MaxCount()
	if gc.segmentRepeatCount > maxCount {
		return s.createError(
			maxSegRepeatCountExceeded,
			fmt.Sprintf("Max. repeat count of segment %s (%d) exceeded: %d",
				segment.Id(), maxCount, gc.segmentRepeatCount))
	} else {
		log.Printf("Repeating segment %s for %dth time", s.currentSegID, gc.segmentRepeatCount)
		return nil
	}
}

func (s *SegSeqValidator) handleRepeatGroup(segment *msg.Seg) error {
	log.Printf("handleRepeatGroup %s", segment)
	gc := s.currentGroupContext()
	if gc.groupRepeatCount >= gc.groupSpecPart.MaxCount() {
		return s.createError(
			maxGroupRepeatCountExceeded,
			fmt.Sprintf("Group segment %s exceeds max group count", segment.Id()))
	} else {
		gc.groupRepeatCount++
		log.Printf("Group repeat count now %d", gc.groupRepeatCount)
		s.nestedMsgBuilder.AddSegGrp(segment.Id())

		//s.incrementCurrentMsgSpecPartIndex()
		return nil
	}
}

func (s *SegSeqValidator) handleSeg(segment *msg.Seg) (matched bool, err error) {
	currentMsgSpecPart := s.getCurrentMsgSpecPart()
	log.Printf("handleSeg %s; current spec: %s",
		segment.Id(), currentMsgSpecPart)
	gc := s.currentGroupContext()
	gc.segmentRepeatCount = 1

	if s.currentSegSpecID != segment.Id() {
		if currentMsgSpecPart.IsMandatory() {
			return false, s.createError(unexpectedSeg,
				fmt.Sprintf("Got unexpected segment %s", segment.Id()))
		} else {
			s.incrementCurrentMsgSpecPartIndex()
			return false, nil
		}
	} else {
		log.Printf("Specs are equal: %s", segment.Id())
	}

	s.nestedMsgBuilder.AddSeg(segment)

	s.setNewState(seqStateSeg)
	return true, nil
}

func (s *SegSeqValidator) getCurrentMsgSpecPart() msp.MsgSpecPart {
	return s.currentGroupContext().currentPart()
}

func (s *SegSeqValidator) nextCurrentMsgSpecPart() msp.MsgSpecPart {
	return s.currentGroupContext().nextPart()
}

func (s *SegSeqValidator) incrementCurrentMsgSpecPartIndex() bool {
	currentGroupContext := s.currentGroupContext()
	if currentGroupContext.AtEnd() {
		log.Printf("Group context %s at end; no increment possible", currentGroupContext)
		return false
	} else {
		currentMsgPart := currentGroupContext.currentPart()
		log.Printf("incr currentMsgSpecPartIndex: %d (%s) --> %d (%s)",
			currentGroupContext.partIndex, currentMsgPart.Name(),
			currentGroupContext.partIndex+1, currentGroupContext.nextPart().Name())
		currentGroupContext.partIndex++
		s.currentSegSpecID = s.getCurrentMsgSpecPart().Id()
		return true
	}
}

func (s *SegSeqValidator) handleStateSeg(
	segID string, segment *msg.Seg,
	msgSpecPart *msp.MsgSpecSegPart) (ret bool, err error) {

	if msgSpecPart.SegSpec.Id == segID {
		if !s.isAtTopLevel() && s.currentGroupContext().groupSpecPart.Id() == segID {
			return true, s.handleRepeatGroup(segment)
		} else {
			s.nestedMsgBuilder.AddSeg(segment)
			return true, s.handleRepSeg(segment)
		}
	} else {
		s.incrementCurrentMsgSpecPartIndex()
		found, err := s.handleSeg(segment)
		if err != nil {
			return true, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func (s *SegSeqValidator) enterGroup(
	msgSpecPart *msp.MsgSpecSegGrpPart) {

	log.Printf("ENTERING GROUP %s", msgSpecPart.Name())
	repeatSegGroup := s.nestedMsgBuilder.AddSegGrp(msgSpecPart.Name())

	groupContext := NewSegSeqGroupContext(
		msgSpecPart, msgSpecPart.Children(),
		repeatSegGroup)
	s.groupStack.Push(groupContext)
}

// Whether validation/construction is currently operation
// at the first level, i.e. not in a group
func (s *SegSeqValidator) isAtTopLevel() bool {
	return s.groupStack.Len() < 2
}

func (s *SegSeqValidator) handleSegGroup(
	segID string,
	msgSpecPart *msp.MsgSpecSegGrpPart) (ret bool, err error) {

	triggerSegId := msgSpecPart.Id()
	if triggerSegId == segID {
		s.enterGroup(msgSpecPart)
	} else {
		if msgSpecPart.IsMandatory() {
			return true, s.createError(
				missingGroup,
				fmt.Sprintf("mandatory group %s missing", triggerSegId))
		} else {
			log.Printf("Skipping group %s (%s)", msgSpecPart.Name(), triggerSegId)
			s.incrementCurrentMsgSpecPartIndex()
		}
	}
	return false, nil
}

func (s *SegSeqValidator) leaveGroup(gc *SegSeqGroupContext) {
	log.Printf("LEAVING GROUP %s", gc.groupSpecPart.Name())
	s.groupStack.Pop()
	s.currentGroupContext().partIndex++
}

func (s *SegSeqValidator) checkGroupStack(segment *msg.Seg) (ret bool) {
	gc := s.currentGroupContext()
	if gc.AtEnd() {
		log.Printf("No more parts in current group spec")
		if !s.isAtTopLevel() {
			if segment.Id() == gc.groupSpecPart.Id() {
				log.Printf("TODO Group repetition")
				s.currentGroupContext().partIndex = 0
			} else {
				s.leaveGroup(gc)
			}

			return false
		} else {
			log.Printf("Returning from top level")
			s.currentGroupContext().partIndex++
			return true
		}
	}
	return false
}

func (s *SegSeqValidator) handleStateSearching(
	segID string, msgSpecPart *msp.MsgSpecSegPart) (ret bool, err error) {

	if msgSpecPart.SegSpec.Id == segID {
		s.setNewState(seqStateSeg)
	} else {
		if msgSpecPart.IsMandatory() {
			return true, s.createError(
				missingMandatorySeg,
				fmt.Sprintf("Mandatory segment %s is missing",
					msgSpecPart.SegSpec.Id))
		}
		s.incrementCurrentMsgSpecPartIndex()
	}
	return false, nil
}

// Advance in spec according to current segment
// Processes a single segment
func (s *SegSeqValidator) processSeg(segment *msg.Seg) error {
	log.Printf("############## processSeg %s", segment.Id())
	log.Printf("\tmessage spec: %s", s.msgSpec)
	log.Printf("\tnested message: %s", s.nestedMsgBuilder.nestedMsg.Dump())
	segID := segment.Id()
	s.currentSegIndex++

	for {
		ret := s.checkGroupStack(segment)
		if ret {
			return nil
		}

		msgSpecPart := s.currentGroupContext().currentPart()
		log.Printf(
			"\tLooping: (state: %s) (stack size: %d, group context: %s) (seg: %s) (%dth); seg spec: %s",
			s.state, s.groupStack.Len(), s.currentGroupContext(), segID,
			s.currentSegIndex, msgSpecPart)
		s.currentSegSpecID = msgSpecPart.Id()

		switch msgSpecPart := msgSpecPart.(type) {
		case *msp.MsgSpecSegPart:
			log.Printf("seg spec ID: %s", s.currentSegSpecID)

			switch s.state {
			case seqStateInitial:
				s.setNewState(seqStateSearching)

			case seqStateGroupStart:
				s.setNewState(seqStateSeg)

			case seqStateSeg:
				ret, err := s.handleStateSeg(segID, segment, msgSpecPart)
				if ret || err != nil {
					return err
				}

			case seqStateSearching:
				ret, err := s.handleStateSearching(segID, msgSpecPart)
				if ret || err != nil {
					return err
				}
			default:
				panic(fmt.Sprintf("Unhandled case: %d", s.state))
			}

		case *msp.MsgSpecSegGrpPart:
			ret, err := s.handleSegGroup(segID, msgSpecPart)
			if ret || err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("Unknown type %T", msgSpecPart))
		}
	}
}

func (s *SegSeqValidator) checkRemainingMandatorySegs() error {
	numSegSpecParts := len(s.msgSpec.TopLevelParts())
	currentMsgSpecPartIndex := s.currentGroupContext().partIndex
	log.Printf("Checking for mandatory segments after message starting at spec index %d",
		currentMsgSpecPartIndex)
	for i := currentMsgSpecPartIndex + 1; i < numSegSpecParts; i++ {
		specPart := s.msgSpec.TopLevelPart(i)
		if specPart.IsMandatory() {
			return s.createError(
				missingMandatorySeg,
				fmt.Sprintf("Mandatory segment %s after end of message",
					specPart.Id()))
		}
	}
	return nil
}

// TODO: return mapping of spec to message segments to allow querying
func (s *SegSeqValidator) Validate(rawMessage *msg.RawMessage) error {
	if len(rawMessage.Segs) == 0 {
		return NewSegSeqError(noSegs, "")
	}
	s.nestedMsgBuilder = NewNestedMsgBuilder(rawMessage.Name, s.groupStack)

	for _, segment := range rawMessage.Segs {
		err := s.processSeg(segment)
		if err != nil {
			return err
		}
	}

	log.Printf("Message ended; TODO check if spec has been fulfilled")
	s.incrementCurrentMsgSpecPartIndex()
	err := s.checkRemainingMandatorySegs()
	if err != nil {
		return err
	}

	log.Printf("Dump of nested message")
	log.Printf(s.nestedMsgBuilder.nestedMsg.Dump())
	log.Printf("%#v", s.nestedMsgBuilder.nestedMsg)
	return nil
}

func NewSegSeqValidator(msgSpec *msp.MsgSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(msgSpec.TopLevelParts()) == 0 {
		return nil, NewSegSeqError(noSegSpecs, "")
	}

	groupContext := NewSegSeqGroupContext(
		msgSpec.TopLevelGroup, msgSpec.TopLevelGroup.Children(), nil)
	groupStack := &util.Stack{}
	groupStack.Push(groupContext)

	return &SegSeqValidator{
		msgSpec:          msgSpec,
		state:            seqStateInitial,
		groupStack:       groupStack,
		nestedMsgBuilder: nil,
	}, nil
}

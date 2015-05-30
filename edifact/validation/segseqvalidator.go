package validation

import (
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	"github.com/bbiskup/edify/edifact/util"
	"log"
)

// Validates segment sequence
// builds structure for navigation/query
type SegSeqValidator struct {
	// index in raw message
	currentSegmentIndex int

	// Specification for message type under validation
	messageSpec *msgspec.MessageSpec

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

func (s *SegSeqValidator) createError(kind SegSeqErrorKind, msg string) error {
	return NewSegSeqError(
		kind, fmt.Sprintf("Error at segment #%d (%s)",
			s.currentSegmentIndex, msg))
}

// Sets new state and logs state transition
func (s *SegSeqValidator) setNewState(newState SegSeqState) {
	log.Printf("State transition %s --> %s", s.state, newState)
	s.state = newState
}

func (s *SegSeqValidator) handleRepeatSegment(segment *msg.Segment) error {
	log.Printf("handleRepeatSegment %s", segment)
	gc := s.currentGroupContext()
	gc.segmentRepeatCount++
	maxCount := s.getCurrentMsgSpecPart().MaxCount()
	if gc.segmentRepeatCount > maxCount {
		return s.createError(
			maxSegmentRepeatCountExceeded,
			fmt.Sprintf("Max. repeat count of segment %s (%d) exceeded: %d",
				segment.Id(), maxCount, gc.segmentRepeatCount))
	} else {
		log.Printf("Repeating segment %s for %dth time", s.currentSegID, gc.segmentRepeatCount)
		return nil
	}
}

func (s *SegSeqValidator) handleRepeatGroup(segment *msg.Segment) error {
	log.Printf("handleRepeatGroup %s", segment)
	gc := s.currentGroupContext()
	if gc.groupRepeatCount >= gc.groupSpecPart.MaxCount() {
		return s.createError(
			maxGroupRepeatCountExceeded,
			fmt.Sprintf("Group segment %s exceeds max group count", segment.Id()))
	} else {
		gc.groupRepeatCount++
		log.Printf("Group repeat count now %d", gc.groupRepeatCount)
		s.nestedMsgBuilder.AddSegment(segment)

		//s.incrementCurrentMsgSpecPartIndex()
		return nil
	}
}

func (s *SegSeqValidator) handleSegment(segment *msg.Segment) (matched bool, err error) {
	currentMsgSpecPart := s.getCurrentMsgSpecPart()
	log.Printf("handleSegment %s; current spec: %s",
		segment.Id(), currentMsgSpecPart)
	gc := s.currentGroupContext()
	gc.segmentRepeatCount = 1

	if s.currentSegSpecID != segment.Id() {
		if currentMsgSpecPart.IsMandatory() {
			return false, s.createError(unexpectedSegment,
				fmt.Sprintf("Got unexpected segment %s", segment.Id()))
		} else {
			s.incrementCurrentMsgSpecPartIndex()
			return false, nil
		}
	} else {
		log.Printf("Specs are equal: %s", segment.Id())
	}

	s.nestedMsgBuilder.AddSegment(segment)

	s.setNewState(seqStateSeg)
	return true, nil
}

func (s *SegSeqValidator) getCurrentMsgSpecPart() msgspec.MessageSpecPart {
	return s.currentGroupContext().currentPart()
}

func (s *SegSeqValidator) nextCurrentMsgSpecPart() msgspec.MessageSpecPart {
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
	segID string, segment *msg.Segment,
	messageSpecPart *msgspec.MessageSpecSegmentPart) (ret bool, err error) {

	if messageSpecPart.SegmentSpec.Id == segID {
		if !s.isAtTopLevel() && s.currentGroupContext().groupSpecPart.Id() == segID {
			return true, s.handleRepeatGroup(segment)
		} else {
			return true, s.handleRepeatSegment(segment)
		}
	} else {
		s.incrementCurrentMsgSpecPartIndex()
		found, err := s.handleSegment(segment)
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
	messageSpecPart *msgspec.MessageSpecSegmentGroupPart) {

	log.Printf("ENTERING GROUP %s", messageSpecPart.Name())
	repeatSegGroup := s.nestedMsgBuilder.AddSegmentGroup(messageSpecPart.Name())

	groupContext := NewSegSeqGroupContext(
		messageSpecPart, messageSpecPart.Children(),
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
	messageSpecPart *msgspec.MessageSpecSegmentGroupPart) (ret bool, err error) {

	triggerSegmentId := messageSpecPart.Id()
	if triggerSegmentId == segID {
		s.enterGroup(messageSpecPart)
	} else {
		if messageSpecPart.IsMandatory() {
			return true, s.createError(
				missingGroup,
				fmt.Sprintf("mandatory group %s missing", triggerSegmentId))
		} else {
			log.Printf("Skipping group %s (%s)", messageSpecPart.Name(), triggerSegmentId)
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

func (s *SegSeqValidator) checkGroupStack(segment *msg.Segment) (ret bool) {
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
	segID string, messageSpecPart *msgspec.MessageSpecSegmentPart) (ret bool, err error) {

	if messageSpecPart.SegmentSpec.Id == segID {
		s.setNewState(seqStateSeg)
	} else {
		if messageSpecPart.IsMandatory() {
			return true, s.createError(
				missingMandatorySegment,
				fmt.Sprintf("Mandatory segment %s is missing",
					messageSpecPart.SegmentSpec.Id))
		}
		s.incrementCurrentMsgSpecPartIndex()
	}
	return false, nil
}

// Advance in spec according to current segment
// Processes a single segment
func (s *SegSeqValidator) processSegment(segment *msg.Segment) error {
	log.Printf("############## processSegment %s", segment.Id())
	log.Printf("\tmessage spec: %s", s.messageSpec)
	log.Printf("\tnested message: %s", s.nestedMsgBuilder.nestedMsg.Dump())
	segID := segment.Id()
	s.currentSegmentIndex++

	for {
		ret := s.checkGroupStack(segment)
		if ret {
			return nil
		}

		messageSpecPart := s.currentGroupContext().currentPart()
		log.Printf(
			"\tLooping: (state: %s) (stack size: %d, group context: %s) (seg: %s) (%dth); seg spec: %s",
			s.state, s.groupStack.Len(), s.currentGroupContext(), segID,
			s.currentSegmentIndex, messageSpecPart)
		s.currentSegSpecID = messageSpecPart.Id()

		switch messageSpecPart := messageSpecPart.(type) {
		case *msgspec.MessageSpecSegmentPart:
			log.Printf("seg spec ID: %s", s.currentSegSpecID)

			switch s.state {
			case seqStateInitial:
				s.setNewState(seqStateSearching)

			case seqStateGroupStart:
				s.setNewState(seqStateSeg)

			case seqStateSeg:
				ret, err := s.handleStateSeg(segID, segment, messageSpecPart)
				if ret || err != nil {
					return err
				}

			case seqStateSearching:
				ret, err := s.handleStateSearching(segID, messageSpecPart)
				if ret || err != nil {
					return err
				}
			default:
				panic(fmt.Sprintf("Unhandled case: %d", s.state))
			}

		case *msgspec.MessageSpecSegmentGroupPart:
			ret, err := s.handleSegGroup(segID, messageSpecPart)
			if ret || err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("Unknown type %T", messageSpecPart))
		}
	}
}

func (s *SegSeqValidator) checkRemainingMandatorySegments() error {
	numSegSpecParts := len(s.messageSpec.Parts)
	currentMsgSpecPartIndex := s.currentGroupContext().partIndex
	log.Printf("Checking for mandatory segments after message starting at spec index %d",
		currentMsgSpecPartIndex)
	for i := currentMsgSpecPartIndex + 1; i < numSegSpecParts; i++ {
		specPart := s.messageSpec.Parts[i]
		if specPart.IsMandatory() {
			return s.createError(
				missingMandatorySegment,
				fmt.Sprintf("Mandatory segment %s after end of message",
					specPart.Id()))
		}
	}
	return nil
}

// TODO: return mapping of spec to message segments to allow querying
func (s *SegSeqValidator) Validate(rawMessage *msg.RawMessage) error {
	if len(rawMessage.Segments) == 0 {
		return NewSegSeqError(noSegments, "")
	}
	s.nestedMsgBuilder = NewNestedMsgBuilder(rawMessage.Name, s.groupStack)

	for _, segment := range rawMessage.Segments {
		err := s.processSegment(segment)
		if err != nil {
			return err
		}
	}

	log.Printf("Message ended; TODO check if spec has been fulfilled")
	s.incrementCurrentMsgSpecPartIndex()
	err := s.checkRemainingMandatorySegments()
	if err != nil {
		return err
	}

	log.Printf("Dump of nested message")
	log.Printf(s.nestedMsgBuilder.nestedMsg.Dump())
	log.Printf("%#v", s.nestedMsgBuilder.nestedMsg)
	return nil
}

func NewSegSeqValidator(messageSpec *msgspec.MessageSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(messageSpec.Parts) == 0 {
		return nil, NewSegSeqError(noSegmentSpecs, "")
	}

	groupContext := NewSegSeqGroupContext(
		nil, messageSpec.Parts, nil)
	groupStack := &util.Stack{}
	groupStack.Push(groupContext)

	return &SegSeqValidator{
		messageSpec:      messageSpec,
		state:            seqStateInitial,
		groupStack:       groupStack,
		nestedMsgBuilder: nil,
	}, nil
}

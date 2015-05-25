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
	currentSegmentIndex int
	messageSpec         *msgspec.MessageSpec
	currentSegSpecID    string
	state               SegSeqState
	currentSegID        string
	groupStack          *util.Stack
}

func (s *SegSeqValidator) currentGroupContext() *SegSeqGroupContext {
	result := s.groupStack.Peek().(*SegSeqGroupContext)
	return result
}

func (s *SegSeqValidator) createError(kind SegSeqErrorKind, msg string) error {
	return NewSegSeqError(
		kind, fmt.Sprintf("Error at segment #%d (%s)",
			s.currentSegmentIndex, msg))
}

func (s *SegSeqValidator) setNewState(newState SegSeqState) {
	log.Printf("State transition %s --> %s", s.state, newState)
	s.state = newState
}

func (s *SegSeqValidator) handleRepeat(segment *msg.Segment) error {
	log.Printf("handleRepeat %s", segment)
	gc := s.currentGroupContext()
	gc.repeatCount++
	maxCount := s.getCurrentMsgSpecPart().MaxCount()
	if gc.repeatCount > maxCount {
		return s.createError(
			maxRepeatCountExceeded,
			fmt.Sprintf("Max. repeat count of segment %s (%d) exceeded: %d",
				s.currentSegID, maxCount, gc.repeatCount))
	} else {
		log.Printf("Repeating segment %s for %dth time", s.currentSegID, gc.repeatCount)
		return nil
	}
}

func (s *SegSeqValidator) handleSegment(segment *msg.Segment) (matched bool, err error) {
	currentMsgSpecPart := s.getCurrentMsgSpecPart()
	log.Printf("handleSegment %s; current spec: %s",
		segment.Id(), currentMsgSpecPart)
	s.currentGroupContext().repeatCount = 1

	if s.currentSegSpecID != segment.Id() {
		if currentMsgSpecPart.IsMandatory() {
			return false, s.createError(unexpectedSegment,
				fmt.Sprintf("Got segment %s", segment.Id()))
		} else {
			s.incrementCurrentMsgSpecPartIndex()
			return false, nil
		}

	}

	s.setNewState(seqStateSeg)
	return true, nil
}

func (s *SegSeqValidator) getCurrentMsgSpecPart() msgspec.MessageSpecPart {
	//return s.messageSpec.Parts[s.currentMsgSpecPartIndex]
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

// Advance in spec according to current segment
// Processes a single segment
func (s *SegSeqValidator) processSegment(segment *msg.Segment) error {
	log.Printf("############## processSegment %s", segment.Id())
	log.Printf("\tmessage spec: %s", s.messageSpec)
	segID := segment.Id()
	s.currentSegmentIndex++

	for {
		if s.currentGroupContext().AtEnd() {
			log.Printf("No more parts in current group spec")
			if s.groupStack.Len() > 1 {
				log.Printf("Leaving group")
				s.groupStack.Pop()
				s.currentGroupContext().partIndex++
				continue
			} else {
				log.Printf("Returning from top level")
				return nil
			}
		}
		messageSpecPart := s.currentGroupContext().currentPart()
		log.Printf("\tLooping: (state: %s) (stack size: %d, group context: %s) (seg: %s) (%dth); seg spec: %s",
			s.state, s.groupStack.Len(), s.currentGroupContext(), segID, s.currentSegmentIndex, messageSpecPart)
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
				if messageSpecPart.SegmentSpec.Id == segID {
					return s.handleRepeat(segment)
				} else {
					s.incrementCurrentMsgSpecPartIndex()
					found, err := s.handleSegment(segment)
					if err != nil {
						return err
					}
					if found {
						return nil
					}
				}

			case seqStateSearching:
				if messageSpecPart.SegmentSpec.Id == segID {
					s.setNewState(seqStateSeg)
				} else {
					if messageSpecPart.IsMandatory() {
						return s.createError(
							missingMandatorySegment,
							fmt.Sprintf("Mandatory segment %s is missing",
								messageSpecPart.SegmentSpec.Id))
					}
					s.incrementCurrentMsgSpecPartIndex()
				}

			default:
				panic(fmt.Sprintf("Unhandled case: %d", s.state))
			}

		case *msgspec.MessageSpecSegmentGroupPart:
			triggerSegmentId := messageSpecPart.TriggerSegmentPart().SegmentSpec.Id
			if triggerSegmentId == segID {
				log.Printf("Entering group %s", messageSpecPart.Name())
				groupContext := &SegSeqGroupContext{
					messageSpecPart, messageSpecPart.Children(), 0, 0}
				s.groupStack.Push(groupContext)
			} else {
				if messageSpecPart.IsMandatory() {
					return s.createError(
						missingGroup,
						fmt.Sprintf("mandatory group %s missing", triggerSegmentId))
				} else {
					log.Printf("Skipping group %s", triggerSegmentId)
					s.incrementCurrentMsgSpecPartIndex()
				}
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
func (s *SegSeqValidator) Validate(message *msg.Message) error {
	if len(message.Segments) == 0 {
		return NewSegSeqError(noSegments, "")
	}
	for _, segment := range message.Segments {
		err := s.processSegment(segment)
		if err != nil {
			return err
		}
	}

	log.Printf("Message ended; TODO check if spec has been fulfilled")
	s.incrementCurrentMsgSpecPartIndex()
	return s.checkRemainingMandatorySegments()
}

func NewSegSeqValidator(messageSpec *msgspec.MessageSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(messageSpec.Parts) == 0 {
		return nil, NewSegSeqError(noSegmentSpecs, "")
	}

	groupContext := &SegSeqGroupContext{nil, messageSpec.Parts, 0, 0}
	groupStack := &util.Stack{}
	groupStack.Push(groupContext)

	return &SegSeqValidator{
		messageSpec: messageSpec,
		state:       seqStateInitial,
		groupStack:  groupStack,
	}, nil
}

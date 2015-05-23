package validation

import (
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	//segspec "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
)

// Validates segment sequence
// builds structure for navigation/query
type SegSeqValidator struct {
	messageSpec             *msgspec.MessageSpec
	message                 *msg.Message
	currentSegmentIndex     int
	currentPartIndex        int // index in current group (or at top level)
	currentMessageSpecParts []msgspec.MessageSpecPart
	repeatCount             int
	previousSegmentId       string
}

type SegSeqErrorKind string

const (
	missingMandatorySegment SegSeqErrorKind = "missing_mandatory_segment"
	noMoreSegments          SegSeqErrorKind = "no_more_segments"
	maxRepeatCountExceeded  SegSeqErrorKind = "max_repeat_count_exceeded"
	missingGroup            SegSeqErrorKind = "missing_group"
	noSegmentSpecs          SegSeqErrorKind = "no_segment_specs"
	noSegments              SegSeqErrorKind = "no_segments"
	//unexpectedSegment       SegSeqErrorKind = "unexpected_segment"
	unexpectedErr SegSeqErrorKind = "unexpected_err"
)

// An exception that provides an error kind to check for specific error conditions
type SegSeqError struct {
	kind    SegSeqErrorKind
	message string
}

func (e SegSeqError) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.message)
}

func NewSegSeqError(kind SegSeqErrorKind, message string) SegSeqError {
	if message == "" {
		message = string(kind)
	}
	return SegSeqError{kind, message}
}

func (s *SegSeqValidator) createError(kind SegSeqErrorKind, msg string) error {
	return NewSegSeqError(kind, fmt.Sprintf("Error at segment #%d (%s)",
		s.currentSegmentIndex, msg))
}

// Searches for given segment ID, and returns the first segment spec part
// with this ID. If a mandatory segment is found before, an error is raised,
// meaning that the segment sequence is incorrect
func (s *SegSeqValidator) advance(segIndex int, segID string) error {
	currentMessageSpecPartsLen := len(s.currentMessageSpecParts)
	log.Printf("advance segIndex = %d, segID = %s, currentMessageSpecPartsLen = %d",
		segIndex, segID, currentMessageSpecPartsLen)
	for i := s.currentPartIndex; i < currentMessageSpecPartsLen; i++ {
		segSpecPart := s.currentMessageSpecParts[i]
		log.Printf("Current segSpecPart: %s", segSpecPart)

		switch segSpecPart := segSpecPart.(type) {
		case *msgspec.MessageSpecSegmentPart:
			log.Printf("At segment spec %s", segSpecPart)
			segSpecID := segSpecPart.SegmentSpec.Id
			log.Printf("segSpecID: %s", segSpecID)
			if segID == s.previousSegmentId {
				// Simple case: repetition
				s.repeatCount++
				log.Printf("Repeating segment type %s (count: %d)",
					segID, s.repeatCount)
				if s.repeatCount > segSpecPart.MaxCount() {
					return s.createError(
						maxRepeatCountExceeded,
						fmt.Sprintf("Max repeat count %d exceeded: %d",
							segSpecPart.MaxCount(), s.repeatCount))
				}
				s.currentPartIndex = i + 1
				return nil
			} else {
				log.Printf("New segment type %s", segID)
				s.repeatCount = 1

				if segSpecID == segID {
					log.Printf("Found segment '%s'", segID)
					s.currentPartIndex = i + 1
					s.previousSegmentId = segID
					return nil
				} else {
					if segSpecPart.IsMandatory() {
						return s.createError(
							missingMandatorySegment,
							fmt.Sprintf("Missing mandatory segment '%s'", segSpecID))
					} else {
						log.Printf("Skipping optional spec segment %s", segSpecID)
					}
				}
			}

		case *msgspec.MessageSpecSegmentGroupPart:
			log.Printf("At group spec %s", segSpecPart)
			triggerSegmentPart := segSpecPart.TriggerSegmentPart()
			log.Printf("### %s: %s", triggerSegmentPart, segSpecPart.Children())
			if segSpecPart.IsMandatory() {
				if triggerSegmentPart.SegmentSpec.Id != segID {
					if triggerSegmentPart.IsMandatory() {
						return s.createError(
							missingGroup,
							fmt.Sprintf("Missing mandatory trigger segment '%s' for group %s",
								triggerSegmentPart, segSpecPart.Name()))
					}
				}
			}
		default:
			panic(fmt.Sprintf("Unsupported spec part type: %T", segSpecPart))

		}
	}
	return s.createError(noMoreSegments, "No more segments")
}

// TODO: return mapping of spec to message segments to allow querying
func (s *SegSeqValidator) Validate(message *msg.Message) error {
	if len(message.Segments) == 0 {
		return NewSegSeqError(noSegments, "")
	}
	for segIndex, segment := range message.Segments {
		s.currentSegmentIndex = segIndex
		segID := segment.Id()
		err := s.advance(segIndex, segID)
		if err != nil {
			return err
		}
	}

	log.Printf("Message ended; checking if spec has been fulfilled")
	return s.advance(s.currentSegmentIndex+1, "___")
}

func NewSegSeqValidator(messageSpec *msgspec.MessageSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(messageSpec.Parts) == 0 {
		return nil, NewSegSeqError(noSegmentSpecs, "")
	}
	return &SegSeqValidator{
		messageSpec:             messageSpec,
		currentMessageSpecParts: messageSpec.Parts,
	}, nil
}

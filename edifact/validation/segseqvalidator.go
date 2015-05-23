package validation

import (
	"errors"
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
			// Simple case: repetition
			segSpecID := segSpecPart.SegmentSpec.Id
			log.Printf("segSpecID: %s", segSpecID)
			if segID == s.previousSegmentId {
				s.repeatCount++
				log.Printf("Repeating segment type %s (count: %d)",
					segID, s.repeatCount)
				if s.repeatCount > segSpecPart.MaxCount() {
					return s.createError(
						maxRepeatCountExceeded,
						fmt.Sprintf("Max repeat count %d exceeded", s.repeatCount))
				}
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
					}
				}
			}

		case *msgspec.MessageSpecSegmentGroupPart:
		default:
			panic(fmt.Sprintf("Unsupported spec part type: %T", segSpecPart))

		}
	}
	return s.createError(noMoreSegments, "No more segments")
}

// TODO: return mapping of spec to message segments to allow querying
func (s *SegSeqValidator) Validate(message *msg.Message) error {
	for segIndex, segment := range message.Segments {
		s.currentSegmentIndex = segIndex
		segID := segment.Id()
		err := s.advance(segIndex, segID)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSegSeqValidator(messageSpec *msgspec.MessageSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(messageSpec.Parts) == 0 {
		return nil, errors.New("No segment specs")
	}
	return &SegSeqValidator{
		messageSpec:             messageSpec,
		currentMessageSpecParts: messageSpec.Parts,
	}, nil
}

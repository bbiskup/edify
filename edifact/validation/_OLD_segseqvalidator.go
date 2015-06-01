package validation

import (
	"fmt"
	msg "github.com/bbiskup/edify/edifact/msg"
	msgspec "github.com/bbiskup/edify/edifact/spec/message"
	segspec "github.com/bbiskup/edify/edifact/spec/segment"
	"log"
)

// Validates segment sequence
// builds structure for navigation/query
type SegSeqValidator struct {
	msgSpec *msgspec.MsgSpec
}

type SegSeqErrorKind string

const (
	missingMandatorySegment SegSeqErrorKind = "missing_mandatory_segment"
	noMoreSegments          SegSeqErrorKind = "no_more_segments"
	maxRepeatCountExceeded  SegSeqErrorKind = "max_repeat_count_exceeded"
	missingGroup            SegSeqErrorKind = "missing_group"
	noSegSpecs          SegSeqErrorKind = "no_segment_specs"
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
		s.currentSegIndex, msg))
}

// TODO: return mapping of spec to message segments to allow querying
func (s *SegSeqValidator) Validate(message *msg.Message) error {
	if len(message.Segments) == 0 {
		return NewSegSeqError(noSegments, "")
	}
	s.message = message
	for _, part := range s.msgSpec.Parts {
		if err != nil {
			return err
		}
	}

	log.Printf("Message ended; TODO check if spec has been fulfilled")
	return nil
}

func NewSegSeqValidator(msgSpec *msgspec.MsgSpec) (segSeqValidator *SegSeqValidator, err error) {
	if len(msgSpec.Parts) == 0 {
		return nil, NewSegSeqError(noSegSpecs, "")
	}
	return &SegSeqValidator{
		msgSpec: msgSpec,
	}, nil
}

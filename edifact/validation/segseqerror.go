package validation

import (
	"fmt"
)

type SegSeqErrKind string

const (
	missingMandatorySegment       SegSeqErrKind = "missing_mandatory_segment"
	maxSegmentRepeatCountExceeded SegSeqErrKind = "max_segment_repeat_count_exceeded"
	maxGroupRepeatCountExceeded   SegSeqErrKind = "max_group_repeat_count_exceeded"
	missingGroup                  SegSeqErrKind = "missing_group"
	noSegSpecs                SegSeqErrKind = "no_segment_specs"
	noSegments                    SegSeqErrKind = "no_segments"
	unexpectedSegment             SegSeqErrKind = "unexpected_segment"
)

// An exception that provides an error kind to check for specific error conditions
type SegSeqError struct {
	kind    SegSeqErrKind
	message string
}

func (e SegSeqError) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.message)
}

func NewSegSeqError(kind SegSeqErrKind, message string) SegSeqError {
	if message == "" {
		message = string(kind)
	}
	return SegSeqError{kind, message}
}

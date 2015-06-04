package validation

import (
	"fmt"
)

type SegSeqErrKind string

const (
	missingMandatorySeg         SegSeqErrKind = "missing_mandatory_segment"
	maxSegRepeatCountExceeded   SegSeqErrKind = "max_segment_repeat_count_exceeded"
	maxGroupRepeatCountExceeded SegSeqErrKind = "max_group_repeat_count_exceeded"
	noSegSpecs                  SegSeqErrKind = "no_segment_specs"
	noSegs                      SegSeqErrKind = "no_segments"
	unexpectedSeg               SegSeqErrKind = "unexpected_segment"
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

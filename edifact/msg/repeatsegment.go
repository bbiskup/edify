package msg

import (
	"bytes"
	"fmt"
)

// A segment that is repeated 1 to n times
type RepeatSegment struct {
	segments []*Segment
}

// From Interface RepeatMsgPart
func (s *RepeatSegment) Count() int {
	return len(s.segments)
}

// From SegmentOrGroup
func (s *RepeatSegment) Id() string {
	return s.segments[0].Id()
}

func (s *RepeatSegment) AddSegment(segment *Segment) {
	s.segments = append(s.segments, segment)
}

func (s *RepeatSegment) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	for repeat, segment := range s.segments {
		buf.WriteString(fmt.Sprintf("%s[%d] %s\n", indentStr, repeat, segment.Id()))
	}
	return buf.String()
}

func NewRepeatSegment(segments ...*Segment) *RepeatSegment {
	return &RepeatSegment{
		segments,
	}
}

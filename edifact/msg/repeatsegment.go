package msg

import (
	"bytes"
	"strconv"
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
		buf.WriteString(indentStr + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
			segment.Id() + "\n")
	}
	return buf.String()
}

func NewRepeatSegment(segments ...*Segment) *RepeatSegment {
	return &RepeatSegment{
		segments,
	}
}

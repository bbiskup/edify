package msg

import (
	"bytes"
	"fmt"
	"github.com/bbiskup/edify/edifact/util"
	"strconv"
)

// A segment that is repeated 1 to n times
type RepSeg struct {
	segments []*Seg
}

func (s *RepSeg) String() string {
	return fmt.Sprintf("RepSeg %s (%dx)", s.Id(), s.Count())
}

// From Interface RepeatMsgPart
func (s *RepSeg) Count() int {
	return len(s.segments)
}

// From SegOrGroup
func (s *RepSeg) Id() string {
	if len(s.segments) == 0 {
		panic("No segments")
	}
	return s.segments[0].Id()
}

// Get n-th repeat
func (s *RepSeg) GetSeg(index int) *Seg {
	return s.segments[index]
}

// Append another repetition of the segment
func (s *RepSeg) AppendSeg(segment *Seg) {
	s.segments = append(s.segments, segment)
}

func (s *RepSeg) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := util.GetIndentStr(indent)
	indentStr2 := util.GetIndentStr(indent + 1)
	buf.WriteString(indentStr + "RepSeg\n")
	for repeat, segment := range s.segments {
		buf.WriteString(indentStr2 + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
			segment.Id() + "\n")
	}
	return buf.String()
}

func NewRepSeg(segments ...*Seg) *RepSeg {
	return &RepSeg{
		segments,
	}
}

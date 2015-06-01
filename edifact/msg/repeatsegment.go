package msg

import (
	"bytes"
	"strconv"
)

// A segment that is repeated 1 to n times
type RepSeg struct {
	segments []*Seg
}

// From Interface RepeatMsgPart
func (s *RepSeg) Count() int {
	return len(s.segments)
}

// From SegOrGroup
func (s *RepSeg) Id() string {
	return s.segments[0].Id()
}

func (s *RepSeg) Get(index int) *Seg {
	return s.segments[index]
}

func (s *RepSeg) AddSeg(segment *Seg) {
	s.segments = append(s.segments, segment)
}

func (s *RepSeg) Dump(indent int) string {
	var buf bytes.Buffer
	indentStr := getIndentStr(indent)
	for repeat, segment := range s.segments {
		buf.WriteString(indentStr + "[" + strconv.FormatInt(int64(repeat), 10) + "] " +
			segment.Id() + "\n")
	}
	return buf.String()
}

func NewRepSeg(segments ...*Seg) *RepSeg {
	return &RepSeg{
		segments,
	}
}

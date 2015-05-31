package util

import (
	"bytes"
	"fmt"
)

// LIFO stack for navigating segment group hierarchy
// Adapted from
// http://www.reddit.com/r/golang/comments/25aeof/building_a_stack_in_go_slices_vs_linked_list/chg21xl
type Stack struct {
	vec []interface{}
}

func (s Stack) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Stack size: %d\n", len(s.vec)))
	for _, item := range s.vec {
		itemStringer, ok := item.(fmt.Stringer)
		var itemStr string
		if ok {
			itemStr = itemStringer.String()
		} else {
			itemStr = fmt.Sprintf("%#v", item)
		}
		buf.WriteString(fmt.Sprintf("\t%s\n", itemStr))
	}
	buf.WriteString("\n")
	return buf.String()
}

func (s Stack) Empty() bool {
	return len(s.vec) == 0
}

func (s Stack) Peek() interface{} {
	return s.vec[len(s.vec)-1]
}

func (s Stack) Len() int {
	return len(s.vec)
}

func (s *Stack) Push(i interface{}) {
	s.vec = append(s.vec, i)
}

func (s *Stack) Pop() interface{} {
	d := s.vec[len(s.vec)-1]
	s.vec = s.vec[:len(s.vec)-1]
	return d
}

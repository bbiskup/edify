package segment

import (
	"fmt"
)

// Segment specification
type SegmentSpec struct {
	Name                string
	Function            string
	DataElementSpecNums []int32
}

func (s *SegmentSpec) String() string {
	return fmt.Sprintf(
		"Segment %s (%d data elems)",
		s.Name, len(s.DataElementSpecNums))
}

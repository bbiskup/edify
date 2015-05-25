package msg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringEmptyMsg(t *testing.T) {
	msg := NewNestedMessage("testname", []SegmentOrGroup{})
	assert.Equal(t, "NestedMessage testname (0 1st-level parts)", msg.String())
}

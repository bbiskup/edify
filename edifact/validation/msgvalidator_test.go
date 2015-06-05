package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMsgTypeFromUNH(t *testing.T) {
	seg := msg.NewSeg("UNH")
	seg.AddElem(msg.NewDataElem([]string{"123"}))
	seg.AddElem(msg.NewDataElem([]string{"ABC", "x", "y"}))
	msgType, err := getMsgTypeFromUNH(seg)
	assert.Nil(t, err)
	assert.Equal(t, "ABC", msgType)
}

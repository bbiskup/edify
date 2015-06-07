package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
)

type SegValidator interface {
	Validate(rawSeg *rawmsg.RawSeg) (*msg.Seg, error)
}

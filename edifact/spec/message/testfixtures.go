package message

import (
	"github.com/bbiskup/edify/edifact/spec/segment"
)

func fixtureSimpleSegmentSpec() *segment.SegmentSpec {
	return segment.NewSegmentSpec("TESTSEGMENT_ID", "TESTSEGMENT_NAME", "tesfunctiont", nil)
}

package message

import (
	"github.com/bbiskup/edify/edifact/segment"
)

func fixtureSimpleSegmentSpec() *segment.SegmentSpec {
	return segment.NewSegmentSpec("TESTSEGMENT_ID", "TESTSEGMENT_NAME", "tesfunctiont", nil)
}

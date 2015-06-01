package message

import (
	"github.com/bbiskup/edify/edifact/spec/segment"
)

func fixtureSimpleSegSpec() *segment.SegSpec {
	return segment.NewSegSpec("TESTSEGMENT_ID", "TESTSEGMENT_NAME", "testfunction", nil)
}

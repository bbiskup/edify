package message

import (
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
)

func fixtureSimpleSegSpec() *ssp.SegSpec {
	return ssp.NewSegSpec("TESTSEGMENT_ID", "TESTSEGMENT_NAME", "testfunction", nil)
}

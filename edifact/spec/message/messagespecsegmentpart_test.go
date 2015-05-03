package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageSpecPartString(t *testing.T) {
	segmentSpec := fixtureSimpleSegmentSpec()
	part := NewMessageSpecSegmentPart(segmentSpec, 5, true, nil)
	assert.Equal(t, "Segment 5 mand. TESTSEGMENT_NAME", part.String())
	assert.Equal(t, 5, part.MaxCount())
	assert.Equal(t, true, part.IsMandatory())
	assert.Equal(t, "TESTSEGMENT_NAME", part.Name())
}

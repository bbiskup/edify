package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/segment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseFile(t *testing.T) {
	segmentSpecs := segment.SegmentSpecMap{} // TODO actual fixture
	parser := NewMessageSpecParser(segmentSpecs)
	spec, err := parser.ParseSpecFile("../../testdata/INVOIC_D.14B")
	assert.Nil(t, err)
	assert.NotNil(t, spec)
	fmt.Printf("Message spec: %s", spec)

	assert.Equal(t, "INVOIC", spec.Id)
	assert.Equal(t, "D", spec.Version)
	assert.Equal(t, "14B", spec.Release)
	assert.Equal(t, "UN", spec.ContrAgency)
	assert.Equal(t, "16", spec.Revision)
	assert.Equal(t, time.Date(2014, time.November, 17, 0, 0, 0, 0, time.UTC), spec.Date)
}

package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryPartStrWithIndex(t *testing.T) {
	p := NewQueryPart(SegmentKind, "abc", 3)
	assert.Equal(t, "QueryPart seg abc 3", p.String())
}

func TestQueryPartStrNoIndex(t *testing.T) {
	p := NewQueryPart(SegmentKind, "abc", noIndex)
	assert.Equal(t, "QueryPart seg abc *", p.String())
}

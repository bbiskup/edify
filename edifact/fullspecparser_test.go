package edifact

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFullSpecParser(t *testing.T) {
	// TODO provide permanent test data
	p, err := NewFullSpecParser("14B", "../testdata/d14b")
	assert.Nil(t, err)

	err = p.Parse()
	assert.Nil(t, err)
	fmt.Printf("p: %#v", p)
}

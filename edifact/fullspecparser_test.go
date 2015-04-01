package edifact

import (
	"fmt"
	"testing"
)

func TestFullSpecParser(t *testing.T) {
	// TODO provide permanent test data
	p, err := NewFullSpecParser("14B", "../testdata/d14b")
	if err != nil {
		t.Errorf("FullSpecParser failed: %s", err)
	}
	fmt.Printf("p: %#v", p)
}

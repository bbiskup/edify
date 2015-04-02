package edifact

import (
	"fmt"
	"testing"
)

func TestFullSpecParser(t *testing.T) {
	// TODO provide permanent test data
	p, err := NewFullSpecParser("14B", "../testdata/d14b")
	if err != nil {
		t.Errorf("NewFullSpecParser failed: %s", err)
	}

	err = p.Parse()
	if err != nil {
		t.Errorf("parse failed: %s", err)
	}
	fmt.Printf("p: %#v", p)
}

package edifact

import (
	"testing"
)

func TestSimpleDataElementString(t *testing.T) {
	elem := NewElement("testName", "testValue")
	res := elem.String()
	expected := "testName testValue"
	if res != expected {
		t.Fatalf("%s != %s")
	}
}

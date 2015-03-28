package tests

import (
	edi "edifice/edifact"
	"testing"
)

func TestString(t *testing.T) {
	elem := edi.NewElement("testName", "testValue")
	res := elem.String()
	expected := "testName testValue"
	if res != expected {
		t.Fatalf("%s != %s")
	}
}

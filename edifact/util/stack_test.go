package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyStack(t *testing.T) {
	s := Stack{}
	assert.Equal(t, 0, s.Len())
	assert.True(t, s.Empty())
}

func TestPush(t *testing.T) {
	s := Stack{}
	val := "val1"
	s.Push(val)
	assert.Equal(t, 1, s.Len())
	assert.False(t, s.Empty())
}

func TestString(t *testing.T) {
	s := Stack{}
	s.Push("one")
	s.Push("two")
	assert.Equal(t, "Stack size: 2\n\t\"one\"\n\t\"two\"\n\n", s.String())
}

func TestStringMixedContent(t *testing.T) {
	s := Stack{}
	s.Push("one")
	s.Push(2)
	assert.Equal(t, "Stack size: 2\n\t\"one\"\n\t2\n\n", s.String())
}

func TestPeek(t *testing.T) {
	s := Stack{}
	s.Push("one")
	assert.Equal(t, "one", s.Peek())
}

func TestPop(t *testing.T) {
	s := Stack{}
	s.Push("one")
	assert.Equal(t, 1, s.Len())
	val := s.Pop()
	assert.Equal(t, "one", val)
	assert.Equal(t, 0, s.Len())
}

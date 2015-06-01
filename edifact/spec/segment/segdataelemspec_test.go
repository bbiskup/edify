package segment

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// implementes dataelement.DataElemSpec
type DummyElem struct {
}

func (d *DummyElem) Id() string {
	return "dummy_id"
}

func (d *DummyElem) Name() string {
	return "dummy_name"
}

func (d *DummyElem) String() string {
	return fmt.Sprintf("%s %s", d.Id(), d.Name())
}

func TestSegDataElemSpec(t *testing.T) {
	e := NewSegDataElemSpec(&DummyElem{}, 3, true)
	const expected = "SegDataElem dummy_id 3x mand."
	assert.Equal(t, expected, e.String())
}

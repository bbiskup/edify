package query

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getNestedMsg(tb testing.TB) *msg.NestedMsg {
	fileName := "../../testdata/messages/INVOIC_1.txt"
	fmt.Printf("EDIFACT file: %s", fileName)
	rawMsg, err := validation.GetRawMsg(fileName)
	require.Nil(tb, err)
	validator := validation.GetTestMsgValidator(tb)
	nestedMsg, err := validator.Validate(rawMsg)
	require.Nil(tb, err)
	return nestedMsg
}

var testNavSpecs = []struct {
	description string
	queryStr    string
	checkFn     func(t *testing.T, msgPart msg.NestedMsgPart, err error)
}{
	{
		"Valid path for segment at top level",
		"seg:BGM[0]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
			assert.Equal(t, "BGM", msgPart.Id())
		},
	},
	{
		"Valid path for segment in group 1",
		"grp:Group_1[0]/seg:RFF[0]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
			assert.Equal(t, "RFF", msgPart.Id())

			seg, ok := msgPart.(*msg.Seg)
			require.True(t, ok)
			cde := seg.DataElems[0].(*msg.CompositeDataElem)
			assert.Equal(t, "C506", cde.Id())

			_1153 := cde.SimpleDataElems[0]
			assert.Equal(t, "1153", _1153.Id())
		},
	},
	{
		"Valid path for composite data element",
		"grp:Group_1[0]/seg:RFF[0]/cmp:C506[0]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
			assert.Equal(t, "C506", msgPart.Id())
			_, ok := msgPart.(*msg.CompositeDataElem)
			assert.True(t, ok)
		},
	},
	{
		"Valid path for simple data element as part of segment",
		"grp:Group_1[0]/seg:RFF[0]/cmp:C506[0]/smp:1153",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
			assert.Equal(t, "1153", msgPart.Id())
			_, ok := msgPart.(*msg.SimpleDataElem)
			assert.True(t, ok)
		},
	},
	{
		"Incorrect segment index",
		"seg:BGM[1]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, msgPart)
			require.NotNil(t, err)
		},
	},
}

func TestNavigatorNavigate(t *testing.T) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(t)

	fmt.Printf("Nested msg: %s", nestedMsg.Dump())

	for _, spec := range testNavSpecs {
		msgPart, err := navigator.Navigate(spec.queryStr, nestedMsg)
		spec.checkFn(t, msgPart, err)
	}
}

func TestNavigatorGetSegTopLevel(t *testing.T) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(t)
	msgPart, err := navigator.GetSeg("seg:BGM[0]", nestedMsg)
	require.Nil(t, err)
	assert.Equal(t, "BGM", msgPart.Id())
}

func TestNavigatorGetSegInGroup(t *testing.T) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(t)
	msgPart, err := navigator.GetSeg("grp:Group_1[0]/seg:RFF[0]", nestedMsg)
	require.Nil(t, err)
	assert.Equal(t, "RFF", msgPart.Id())
}

func TestNavigatorGetSegGroup(t *testing.T) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(t)
	msgPart, err := navigator.GetSegGrp("grp:Group_1[0]", nestedMsg)
	require.Nil(t, err)
	assert.Equal(t, "Group_1", msgPart.Id())
}

func BenchmarkNavigate1(b *testing.B) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgPart, err := navigator.GetSeg("grp:Group_1[0]/seg:RFF[0]", nestedMsg)
		require.NotNil(b, msgPart)
		require.Nil(b, err)
	}
}

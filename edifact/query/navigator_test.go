package query

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getNestedMsg(t *testing.T) *msg.NestedMsg {
	fileName := "../../testdata/messages/INVOIC_1.txt"
	fmt.Printf("EDIFACT file: %s", fileName)
	rawMsg, err := validation.GetRawMsg(fileName)
	require.Nil(t, err)
	validator := validation.GetTestMsgValidator(t)
	nestedMsg, err := validator.Validate(rawMsg)
	require.Nil(t, err)
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
		"grp:Group_1[0]|seg:RFF[0]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
			assert.Equal(t, "RFF", msgPart.Id())

			seg, ok := msgPart.(*msg.Seg)
			assert.True(t, ok)
			cde := seg.DataElems[0].(*msg.CompositeDataElem)
			assert.Equal(t, "C506", cde.Id())

			_1153 := cde.SimpleDataElems[0]
			assert.Equal(t, "1153", _1153.Id())
		},
	},
	// {
	// 	"Valid path for segment in group 1",
	// 	"grp:Group_1[0]|seg:RFF[0]|cmp:C506[0]",
	// 	func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
	// 		require.Nil(t, err)
	// 		require.NotNil(t, msgPart)
	// 		assert.Equal(t, "C506", msgPart.Id())
	// 	},
	// },
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
		msgPart, err := navigator.navigate(spec.queryStr, nestedMsg)
		spec.checkFn(t, msgPart, err)
	}
}

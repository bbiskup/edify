package query

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/validation"
	"github.com/stretchr/testify/require"
	"testing"
)

func getNestedMsg(t *testing.T) *msg.NestedMsg {
	fileName := "../../testdata/messages/INVOIC_1.txt"
	fmt.Printf("EDIFACT file: %s", fileName)
	rawMsg, err := validation.GetRawMsg(fileName)
	require.Nil(t, err)
	validator := validation.GetValidator(t)
	nestedMsg, err := validator.Validate(rawMsg)
	require.Nil(t, err)
	return nestedMsg
}

var testNavSpecs = []struct {
	queryStr string
	checkFn  func(t *testing.T, msgPart msg.NestedMsgPart, err error)
}{
	{"seg:BGM[0]",
		func(t *testing.T, msgPart msg.NestedMsgPart, err error) {
			require.Nil(t, err)
			require.NotNil(t, msgPart)
		},
	},
}

func TestNavigator(t *testing.T) {
	navigator := NewNavigator()
	nestedMsg := getNestedMsg(t)

	for _, spec := range testNavSpecs {
		msgPart, err := navigator.GetSeg("seg:BGM[0]", nestedMsg)
		spec.checkFn(t, msgPart, err)
	}
}

// package validation

// import (
// 	"github.com/bbiskup/edify/edifact/msg"
// 	"github.com/bbiskup/edify/edifact/spec/codes"
// 	de "github.com/bbiskup/edify/edifact/spec/dataelement"
// 	"github.com/bbiskup/edify/edifact/spec/segment"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"testing"
// )

// func TestValidateSegment(t *testing.T) {
// 	// message
// 	segment1 := msg.NewSegment("ABC")
// 	segment1.AddElement(
// 		msg.NewDataElement([]string{
// 			"abc",
// 			"1",
// 		}),
// 	)

// 	// spec
// 	de1Spec := codes.NewCodesSpec("100", "testcode_1", "testcode_1_desc",
// 		[]*codes.CodeSpec{
// 			codes.NewCodeSpec("1", "value_1", "descr_1"),
// 			codes.NewCodeSpec("2", "value_2", "descr_2"),
// 		})
// 	/*csMap := codes.CodesSpecMap{
// 		"100": de1Spec,
// 	}*/

// 	de0, err := de.NewSimpleDataElementSpec(
// 		"simple_1", "simple_1_name", "simple_1_descr", de.NewRepr(de.Alpha, true, 10), nil)
// 	require.Nil(t, err)

// 	de1, err := de.NewSimpleDataElementSpec(
// 		"simple_2", "simple_2_name", "simple_2_descr", de.NewRepr(de.AlphaNum, true, 1), de1Spec)
// 	require.Nil(t, err)

// 	segDataElemSpecs := []*segment.SegmentDataElementSpec{
// 		segment.NewSegmentDataElementSpec(de0, 1, true),
// 		segment.NewSegmentDataElementSpec(de1, 1, true),
// 	}

// 	segSpec := segment.NewSegmentSpec("ABC", "ABC_segment", "abc_function", segDataElemSpecs)
// 	segSpecMap := segment.SegmentSpecMap{"ABC": segSpec}

// 	validator := NewSegmentValidator(segSpecMap)

// 	valid, err := validator.Validate(segment1)
// 	assert.True(t, valid)
// 	assert.Nil(t, err)
// }

package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/spec/codes"
	de "github.com/bbiskup/edify/edifact/spec/dataelement"
	"github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// fixture
func getValidSegment(t testing.TB) *msg.Segment {
	seg := msg.NewSegment("ABC")
	seg.AddElement(
		msg.NewDataElement([]string{
			"abc",
		}),
	)

	seg.AddElement(
		msg.NewDataElement([]string{
			"1",
		}),
	)
	return seg
}

// fixture: non-existant code
func getInvalidSegmentNonExistantCode(t testing.TB) *msg.Segment {
	seg := msg.NewSegment("ABC")
	seg.AddElement(
		msg.NewDataElement([]string{
			"abc",
		}),
	)

	seg.AddElement(
		msg.NewDataElement([]string{
			"3", // does not exist
		}),
	)
	return seg
}

// fixture: non-existant code
func getInvalidSegmentIncorrectRepr(t testing.TB) *msg.Segment {
	seg := msg.NewSegment("ABC")
	seg.AddElement(
		msg.NewDataElement([]string{
			"abc",
		}),
	)

	seg.AddElement(
		msg.NewDataElement([]string{
			"x", // should be numeric
		}),
	)
	return seg
}

func getSegmentSpecMap(t testing.TB) segment.SegmentSpecMap {
	de1Spec := codes.NewCodesSpec("100", "testcode_1", "testcode_1_desc",
		[]*codes.CodeSpec{
			codes.NewCodeSpec("1", "value_1", "descr_1"),
			codes.NewCodeSpec("2", "value_2", "descr_2"),
		})

	de0, err := de.NewSimpleDataElementSpec(
		"simple_1", "simple_1_name", "simple_1_descr", de.NewRepr(de.Alpha, true, 10), nil)
	require.Nil(t, err)

	de1, err := de.NewSimpleDataElementSpec(
		"simple_2", "simple_2_name", "simple_2_descr", de.NewRepr(de.Num, true, 1), de1Spec)
	require.Nil(t, err)

	segDataElemSpecs := []*segment.SegmentDataElementSpec{
		segment.NewSegmentDataElementSpec(de0, 1, true),
		segment.NewSegmentDataElementSpec(de1, 1, true),
	}

	segSpec := segment.NewSegmentSpec("ABC", "ABC_segment", "abc_function", segDataElemSpecs)
	return segment.SegmentSpecMap{"ABC": segSpec}
}

func TestValidateValidSegment(t *testing.T) {
	segSpecMap := getSegmentSpecMap(t)
	segment := getValidSegment(t)
	validator := NewSegmentValidator(segSpecMap)
	err := validator.Validate(segment)
	assert.Nil(t, err)
}

func TestValidateInvalidSegmentNonExistantCode(t *testing.T) {
	segSpecMap := getSegmentSpecMap(t)
	segment := getInvalidSegmentNonExistantCode(t)
	validator := NewSegmentValidator(segSpecMap)
	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func TestValidateInvalidSegmentIncorrectRepr(t *testing.T) {
	segSpecMap := getSegmentSpecMap(t)
	segment := getInvalidSegmentIncorrectRepr(t)
	validator := NewSegmentValidator(segSpecMap)
	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func BenchmarkParseValidSegment(b *testing.B) {
	segSpecMap := getSegmentSpecMap(b)
	segment := getInvalidSegmentNonExistantCode(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validator := NewSegmentValidator(segSpecMap)
		err := validator.Validate(segment)
		assert.NotNil(b, err)
	}
}

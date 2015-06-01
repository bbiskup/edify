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
func getValidSeg(t testing.TB) *msg.Seg {
	seg := msg.NewSeg("ABC")
	seg.AddElem(
		msg.NewDataElem([]string{
			"abc",
		}),
	)

	seg.AddElem(
		msg.NewDataElem([]string{
			"1",
		}),
	)
	return seg
}

// fixture: non-existant code
func getInvalidSegNonExistantCode(t testing.TB) *msg.Seg {
	seg := msg.NewSeg("ABC")
	seg.AddElem(
		msg.NewDataElem([]string{
			"abc",
		}),
	)

	seg.AddElem(
		msg.NewDataElem([]string{
			"3", // does not exist
		}),
	)
	return seg
}

// fixture: non-existant code
func getInvalidSegIncorrectRepr(t testing.TB) *msg.Seg {
	seg := msg.NewSeg("ABC")
	seg.AddElem(
		msg.NewDataElem([]string{
			"abc",
		}),
	)

	seg.AddElem(
		msg.NewDataElem([]string{
			"x", // should be numeric
		}),
	)
	return seg
}

func getSegSpecMap(t testing.TB) segment.SegSpecMap {
	de1Spec := codes.NewCodesSpec("100", "testcode_1", "testcode_1_desc",
		[]*codes.CodeSpec{
			codes.NewCodeSpec("1", "value_1", "descr_1"),
			codes.NewCodeSpec("2", "value_2", "descr_2"),
		})

	de0, err := de.NewSimpleDataElemSpec(
		"simple_1", "simple_1_name", "simple_1_descr", de.NewRepr(de.Alpha, true, 10), nil)
	require.Nil(t, err)

	de1, err := de.NewSimpleDataElemSpec(
		"simple_2", "simple_2_name", "simple_2_descr", de.NewRepr(de.Num, true, 1), de1Spec)
	require.Nil(t, err)

	segDataElemSpecs := []*segment.SegDataElemSpec{
		segment.NewSegDataElemSpec(de0, 1, true),
		segment.NewSegDataElemSpec(de1, 1, true),
	}

	segSpec := segment.NewSegSpec("ABC", "ABC_segment", "abc_function", segDataElemSpecs)
	return segment.SegSpecMap{"ABC": segSpec}
}

func TestValidateValidSeg(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	segment := getValidSeg(t)
	validator := NewSegValidatorImpl(segSpecMap)
	err := validator.Validate(segment)
	assert.Nil(t, err)
}

func TestValidateInvalidSegNonExistantCode(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	segment := getInvalidSegNonExistantCode(t)
	validator := NewSegValidatorImpl(segSpecMap)
	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func TestValidateInvalidSegIncorrectRepr(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	segment := getInvalidSegIncorrectRepr(t)
	validator := NewSegValidatorImpl(segSpecMap)
	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func BenchmarkParseValidSeg(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	segment := getInvalidSegNonExistantCode(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validator := NewSegValidatorImpl(segSpecMap)
		err := validator.Validate(segment)
		assert.NotNil(b, err)
	}
}

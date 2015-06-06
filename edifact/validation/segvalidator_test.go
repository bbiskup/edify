package validation

import (
	"github.com/bbiskup/edify/edifact/msg"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
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

func getSegSpecMap(t testing.TB) ssp.SegSpecMap {
	de1Spec := csp.NewCodesSpec("100", "testcode_1", "testcode_1_desc",
		[]*csp.CodeSpec{
			csp.NewCodeSpec("1", "value_1", "descr_1"),
			csp.NewCodeSpec("2", "value_2", "descr_2"),
		})

	de0, err := dsp.NewSimpleDataElemSpec(
		"simple_1", "simple_1_name", "simple_1_descr", dsp.NewRepr(dsp.Alpha, true, 10), nil)
	require.Nil(t, err)

	de1, err := dsp.NewSimpleDataElemSpec(
		"simple_2", "simple_2_name", "simple_2_descr", dsp.NewRepr(dsp.Num, true, 1), de1Spec)
	require.Nil(t, err)

	segDataElemSpecs := []*ssp.SegDataElemSpec{
		ssp.NewSegDataElemSpec(de0, 1, true),
		ssp.NewSegDataElemSpec(de1, 1, true),
	}

	segSpec := ssp.NewSegSpec("ABC", "ABC_segment", "abc_function", segDataElemSpecs)
	return ssp.SegSpecMap{"ABC": segSpec}
}

func TestValidateValidSeg(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	segment := getValidSeg(t)

	err := validator.Validate(segment)
	assert.Nil(t, err)
}

func TestValidateInvalidSegNonExistantCode(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	segment := getInvalidSegNonExistantCode(t)

	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func TestValidateInvalidSegIncorrectRepr(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	segment := getInvalidSegIncorrectRepr(t)
	err := validator.Validate(segment)
	assert.NotNil(t, err)
}

func BenchmarkParseInvalidSeg(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	segment := getInvalidSegNonExistantCode(b)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := validator.Validate(segment)
		assert.NotNil(b, err)
	}
}

func BenchmarkParseValidSeg(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	segment := getValidSeg(b)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	assert.NotNil(b, validator)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := validator.Validate(segment)
		assert.Nil(b, err)
	}
}

func BenchmarkNewSegValidatorImpl(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
		assert.NotNil(b, validator)
	}
}

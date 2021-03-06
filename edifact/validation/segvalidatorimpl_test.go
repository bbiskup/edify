package validation

import (
	//"github.com/bbiskup/edify/edifact/msg"
	"github.com/bbiskup/edify/edifact/rawmsg"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
	dsp "github.com/bbiskup/edify/edifact/spec/dataelement"
	ssp "github.com/bbiskup/edify/edifact/spec/segment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// fixture
func getValidRawSeg(t testing.TB) *rawmsg.RawSeg {
	rawSeg := rawmsg.NewRawSeg("ABC")
	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"abc",
		}),
	)

	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"1",
		}),
	)
	return rawSeg
}

// fixture: non-existant code
func getInvalidSegNonExistantCode(t testing.TB) *rawmsg.RawSeg {
	rawSeg := rawmsg.NewRawSeg("ABC")
	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"abc",
		}),
	)

	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"3", // does not exist
		}),
	)
	return rawSeg
}

// fixture: non-existant code
func getInvalidSegIncorrectRepr(t testing.TB) *rawmsg.RawSeg {
	rawSeg := rawmsg.NewRawSeg("ABC")
	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"abc",
		}),
	)

	rawSeg.AddElem(
		rawmsg.NewRawDataElem([]string{
			"x", // should be numeric
		}),
	)
	return rawSeg
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
	rawSeg := getValidRawSeg(t)

	seg, err := validator.Validate(rawSeg)
	assert.Nil(t, err)
	assert.NotNil(t, seg)
}

func TestValidateInvalidSegNonExistantCode(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	rawSeg := getInvalidSegNonExistantCode(t)

	seg, err := validator.Validate(rawSeg)
	assert.Nil(t, seg)
	assert.NotNil(t, err)
}

func TestValidateInvalidSegIncorrectRepr(t *testing.T) {
	segSpecMap := getSegSpecMap(t)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	rawSeg := getInvalidSegIncorrectRepr(t)
	seg, err := validator.Validate(rawSeg)
	assert.Nil(t, seg)
	assert.NotNil(t, err)
}

func BenchmarkParseInvalidSeg(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	rawSeg := getInvalidSegNonExistantCode(b)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		seg, err := validator.Validate(rawSeg)
		assert.Nil(b, seg)
		assert.NotNil(b, err)
	}
}

func BenchmarkParseValidSeg(b *testing.B) {
	segSpecMap := getSegSpecMap(b)
	rawSeg := getValidRawSeg(b)
	validator := NewSegValidatorImpl(ssp.NewSegSpecProviderImpl(segSpecMap))
	assert.NotNil(b, validator)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		seg, err := validator.Validate(rawSeg)
		assert.NotNil(b, seg)
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

package validation

/*
import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageValidatorAUTHOR(t *testing.T) {
	validator, err := NewMessageValidator(getMessageSpec("AUTHOR_D.14B"))
	require.Nil(t, err)
	require.NotNil(t, validator)

	//expected := "^(UNH:)(BGM:)(DTM:){0,1}(BUS:){0,1}((RFF:)(DTM:){0,1}){0,2}((FII:)(CTA:){0,1}(COM:){0,5}){0,5}((NAD:)(CTA:){0,1}(COM:){0,5}){0,3}((LIN:)((RFF:)(DTM:){0,1}){0,5}((SEQ:)(GEI:)(DTM:){0,2}(MOA:){0,1}(DOC:){0,5}){0,}((FII:)(CTA:){0,1}(COM:){0,5}){0,2}((NAD:)(CTA:){0,1}(COM:){0,5}){0,2}){1,}(CNT:){0,5}((AUT:)(DTM:){0,1}){0,5}(UNT:)$"
	//expected := "^(UNH:)(BGM:)(DTM:)*(BUS:)*((RFF:)(DTM:)*){0,2}((FII:)(CTA:)*(COM:){0,5}){0,5}((NAD:)(CTA:)*(COM:){0,5}){0,3}((LIN:)((RFF:)(DTM:)*){0,5}((SEQ:)(GEI:)(DTM:){0,2}(MOA:)*(DOC:){0,5})*((FII:)(CTA:)*(COM:){0,5}){0,2}((NAD:)(CTA:)*(COM:){0,5}){0,2})+(CNT:){0,5}((AUT:)(DTM:)*){0,5}(UNT:)$"
	//expected := "^(UNH:)(BGM:)(DTM:)*(BUS:)*(?P<Group_1>(RFF:)(DTM:)*){0,2}(?P<Group_2>(FII:)(CTA:)*(COM:){0,5}){0,5}(?P<Group_3>(NAD:)(CTA:)*(COM:){0,5}){0,3}(?P<Group_4>(LIN:)(?P<Group_5>(RFF:)(DTM:)*){0,5}(?P<Group_6>(SEQ:)(GEI:)(DTM:){0,2}(MOA:)*(DOC:){0,5})*(?P<Group_7>(FII:)(CTA:)*(COM:){0,5}){0,2}(?P<Group_8>(NAD:)(CTA:)*(COM:){0,5}){0,2})+(CNT:){0,5}(?P<Group_9>(AUT:)(DTM:)*){0,5}(UNT:)$"
	//expected := `^(UNH-[0-9]+:)(BGM-[0-9]+:)(DTM-[0-9]+:)*(BUS-[0-9]+:)*(?P<Group_1>(RFF-[0-9]+:)(DTM-[0-9]+:)*){0,2}(?P<Group_2>(FII-[0-9]+:)(CTA-[0-9]+:)*(COM-[0-9]+:){0,5}){0,5}(?P<Group_3>(NAD-[0-9]+:)(CTA-[0-9]+:)*(COM-[0-9]+:){0,5}){0,3}(?P<Group_4>(LIN-[0-9]+:)(?P<Group_5>(RFF-[0-9]+:)(DTM-[0-9]+:)*){0,5}(?P<Group_6>(SEQ-[0-9]+:)(GEI-[0-9]+:)(DTM-[0-9]+:){0,2}(MOA-[0-9]+:)*(DOC-[0-9]+:){0,5})*(?P<Group_7>(FII-[0-9]+:)(CTA-[0-9]+:)*(COM-[0-9]+:){0,5}){0,2}(?P<Group_8>(NAD-[0-9]+:)(CTA-[0-9]+:)*(COM-[0-9]+:){0,5}){0,2})+(CNT-[0-9]+:){0,5}(?P<Group_9>(AUT-[0-9]+:)(DTM-[0-9]+:)*){0,5}(UNT-[0-9]+:)$`
	//expected := `^(?P<toplevel__UNH>UNH-[0-9]+:)(?P<toplevel__BGM>BGM-[0-9]+:)(?P<toplevel__DTM>DTM-[0-9]+:)*(?P<toplevel__BUS>BUS-[0-9]+:)*((?P<Group_1__RFF>RFF-[0-9]+:)(?P<Group_1__DTM>DTM-[0-9]+:)*){0,2}((?P<Group_2__FII>FII-[0-9]+:)(?P<Group_2__CTA>CTA-[0-9]+:)*(?P<Group_2__COM>COM-[0-9]+:){0,5}){0,5}((?P<Group_3__NAD>NAD-[0-9]+:)(?P<Group_3__CTA>CTA-[0-9]+:)*(?P<Group_3__COM>COM-[0-9]+:){0,5}){0,3}((?P<Group_4__LIN>LIN-[0-9]+:)((?P<Group_5__RFF>RFF-[0-9]+:)(?P<Group_5__DTM>DTM-[0-9]+:)*){0,5}((?P<Group_6__SEQ>SEQ-[0-9]+:)(?P<Group_6__GEI>GEI-[0-9]+:)(?P<Group_6__DTM>DTM-[0-9]+:){0,2}(?P<Group_6__MOA>MOA-[0-9]+:)*(?P<Group_6__DOC>DOC-[0-9]+:){0,5})*((?P<Group_7__FII>FII-[0-9]+:)(?P<Group_7__CTA>CTA-[0-9]+:)*(?P<Group_7__COM>COM-[0-9]+:){0,5}){0,2}((?P<Group_8__NAD>NAD-[0-9]+:)(?P<Group_8__CTA>CTA-[0-9]+:)*(?P<Group_8__COM>COM-[0-9]+:){0,5}){0,2})+(?P<toplevel__CNT>CNT-[0-9]+:){0,5}((?P<Group_9__AUT>AUT-[0-9]+:)(?P<Group_9__DTM>DTM-[0-9]+:)*){0,5}(?P<toplevel__UNT>UNT-[0-9]+:)$`
	//expected := `^(?P<toplevel>UNH-[0-9]+:)(?P<toplevel>BGM-[0-9]+:)(?P<toplevel>DTM-[0-9]+:)*(?P<toplevel>BUS-[0-9]+:)*((?P<Group_1>RFF-[0-9]+:)(?P<Group_1>DTM-[0-9]+:)*){0,2}((?P<Group_2>FII-[0-9]+:)(?P<Group_2>CTA-[0-9]+:)*(?P<Group_2>COM-[0-9]+:)*)*((?P<Group_3>NAD-[0-9]+:)(?P<Group_3>CTA-[0-9]+:)*(?P<Group_3>COM-[0-9]+:)*){0,3}((?P<Group_4>LIN-[0-9]+:)((?P<Group_5>RFF-[0-9]+:)(?P<Group_5>DTM-[0-9]+:)*)*((?P<Group_6>SEQ-[0-9]+:)(?P<Group_6>GEI-[0-9]+:)(?P<Group_6>DTM-[0-9]+:){0,2}(?P<Group_6>MOA-[0-9]+:)*(?P<Group_6>DOC-[0-9]+:)*)*((?P<Group_7>FII-[0-9]+:)(?P<Group_7>CTA-[0-9]+:)*(?P<Group_7>COM-[0-9]+:)*){0,2}((?P<Group_8>NAD-[0-9]+:)(?P<Group_8>CTA-[0-9]+:)*(?P<Group_8>COM-[0-9]+:)*){0,2})+(?P<toplevel>CNT-[0-9]+:)*((?P<Group_9>AUT-[0-9]+:)(?P<Group_9>DTM-[0-9]+:)*)*(?P<toplevel>UNT-[0-9]+:)$`
	expected := `^(?P<toplevel>(?:UNH-[0-9]+:))(?P<toplevel>(?:BGM-[0-9]+:))(?P<toplevel>(?:DTM-[0-9]+:)*)(?P<toplevel>(?:BUS-[0-9]+:)*)((?:(?P<Group_1>(?:RFF-[0-9]+:))(?P<Group_1>(?:DTM-[0-9]+:)*)){0,2})((?:(?P<Group_2>(?:FII-[0-9]+:))(?P<Group_2>(?:CTA-[0-9]+:)*)(?P<Group_2>(?:COM-[0-9]+:)*))*)((?:(?P<Group_3>(?:NAD-[0-9]+:))(?P<Group_3>(?:CTA-[0-9]+:)*)(?P<Group_3>(?:COM-[0-9]+:)*)){0,3})((?:(?P<Group_4>(?:LIN-[0-9]+:))((?:(?P<Group_5>(?:RFF-[0-9]+:))(?P<Group_5>(?:DTM-[0-9]+:)*))*)((?:(?P<Group_6>(?:SEQ-[0-9]+:))(?P<Group_6>(?:GEI-[0-9]+:))(?P<Group_6>(?:DTM-[0-9]+:){0,2})(?P<Group_6>(?:MOA-[0-9]+:)*)(?P<Group_6>(?:DOC-[0-9]+:)*))*)((?:(?P<Group_7>(?:FII-[0-9]+:))(?P<Group_7>(?:CTA-[0-9]+:)*)(?P<Group_7>(?:COM-[0-9]+:)*)){0,2})((?:(?P<Group_8>(?:NAD-[0-9]+:))(?P<Group_8>(?:CTA-[0-9]+:)*)(?P<Group_8>(?:COM-[0-9]+:)*)){0,2}))+)(?P<toplevel>(?:CNT-[0-9]+:)*)((?:(?P<Group_9>(?:AUT-[0-9]+:))(?P<Group_9>(?:DTM-[0-9]+:)*))*)(?P<toplevel>(?:UNT-[0-9]+:))$`
	assert.Equal(t, expected, validator.segmentValidationRegexpStr)
}

func TestMessageValidatorGOVCBR(t *testing.T) {
	validator, err := NewMessageValidator(getMessageSpec("GOVCBR_D.14B"))
	require.Nil(t, err)
	require.NotNil(t, validator)
}

var segListStrSpec = []struct {
	segmentIDs []string
	expected   string
}{
	{[]string{}, ""},
	{[]string{"AAA"}, "AAA-0:"},
	{[]string{"AAA", "BBB"}, "AAA-0:BBB-1:"},
}

func TestBuildSegmentListStr(t *testing.T) {
	for _, spec := range segListStrSpec {
		result, indexSegmentMap := buildSegmentListStr(mapToSegments(spec.segmentIDs))
		assert.Equal(t, spec.expected, result)
		assert.Equal(t, len(spec.segmentIDs), len(indexSegmentMap))
	}
}

var authorSegSeqSpec_OLD = []struct {
	descr       string
	segmentIDs  []string
	valid       bool
	expectError bool
}{
	{"minimal message (only mandatory segments)",
		[]string{
			"UNH", "BGM",
			// Group 1
			"LIN",
			"UNT",
		}, true, false},

	{
		"Mostly mandatory",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS", // both conditional
			// Group 4
			"LIN",
			"UNT",
		}, true, false},

	{
		"Mostly mandatory; one conditional group",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 1
			"LIN",
			// Group 2
			"FII", "CTA", "COM",

			"UNT",
		}, true, false},

	{
		"Some repeat counts > 1",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 4
			"LIN", "LIN", "LIN", "LIN",
			// Group 7
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",

			"UNT",
		}, true, false},

	{
		"No segments at all",
		[]string{}, false, true},

	{
		"Missing mandatory segments",
		[]string{"UNH"}, false, false},

	{
		"First mandatory segment repeated too often",
		[]string{
			"UNH", "UNH", "BGM",
			// Group 1
			"LIN",
			"UNT",
		}, false, false},

	{
		"group 7 repeated too often",
		[]string{
			"UNH", "BGM",
			"DTM", "BUS",
			// Group 4
			"LIN", "LIN", "LIN", "LIN",
			// Group 7
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",
			"FII", "CTA", "COM", "COM", "COM",

			"UNT",
		}, false, false},
}

func TestValidateSegmentList(t *testing.T) {
	validator, err := NewMessageValidator(getMessageSpec("AUTHOR_D.14B"))
	require.Nil(t, err)
	// fmt.Printf("regexp str %s", validator.segmentValidationRegexpStr)
	for _, spec := range authorSegSeqSpec_OLD {
		fmt.Printf("spec: %#v\n", spec)
		result, err := validator.ValidateSegmentList(mapToSegments(spec.segmentIDs))
		assert.Equal(t, spec.expectError, err != nil)
		assert.Equal(t, spec.valid, result)
	}
}

// Benchmark creation of validation regexp for smaller message spec
func BenchmarkNewMessageValidatorSmallSpec(b *testing.B) {
	messageSpec := getMessageSpec("AUTHOR_D.14B")
	for i := 0; i < b.N; i++ {
		validator, err := NewMessageValidator(messageSpec)
		require.Nil(b, err)
		require.NotNil(b, validator)
	}
}

// Benchmark creation of validation regexp for larger UNCE message spec
func BenchmarkNewMessageValidatorLargeSpec(b *testing.B) {
	messageSpec := getMessageSpec("PRICAT_D.14B")
	for i := 0; i < b.N; i++ {
		validator, err := NewMessageValidator(messageSpec)
		require.Nil(b, err)
		require.NotNil(b, validator)
	}
}

// Does not work; invalid repeat count when parsing regexp
// 	same holds for (descending size order)
// 	IFCSUM, ORDRSP, ORDERS, ORDCHG,
//
// // Benchmark creation of validation regexp for largest UNCE message spec
// func BenchmarkNewMessageValidatorLargestSpec(b *testing.B) {
// 	messageSpec := getMessageSpec("GOVCBR_D.14B")
// 	for i := 0; i < b.N; i++ {
// 		validator, err := NewMessageValidator(messageSpec)
// 		require.Nil(b, err)
// 		require.NotNil(b, validator)
// 	}
// }

func BenchmarkValidateAuthorSegments(b *testing.B) {
	validator, err := NewMessageValidator(getMessageSpec("AUTHOR_D.14B"))
	require.Nil(b, err)
	segmentIDs := []string{
		"UNH", "BGM",
		"DTM", "BUS",
		// Group 4
		"LIN", "LIN", "LIN", "LIN",
		// Group 7
		"FII", "CTA", "COM", "COM", "COM",
		"FII", "CTA", "COM", "COM", "COM",

		"UNT"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		valid, err := validator.ValidateSegmentList(mapToSegments(segmentIDs))
		if i == 0 {
			assert.True(b, valid)
			assert.Nil(b, err)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

var getRegexpRepeatStrSpec = []struct {
	minSpecRepeat int
	maxSpecRepeat int
	isGroup       bool
	expected      string
	expectError   bool
}{
	// Segment specs
	{0, 1, false, "*", false},
	{1, 1, false, "", false},
	{0, 2, false, "{0,2}", false},
	{1, 2, false, "{1,2}", false},
	{1, 99, false, "+", false},
	{0, 99, false, "*", false},
	{2, 1, false, "", true},   // max > min
	{1, 100, false, "", true}, // max > allowed max 99 for segment specs

	// Group specs
	{0, 1, false, "*", false},
	{1, 1, false, "", false},
	{0, 2, false, "{0,2}", false},
	{1, 2, false, "{1,2}", false},
	{1, 999, false, "+", false},
	{0, 999, false, "*", false},
	{2, 1, false, "", true},    // max > min
	{1, 1000, false, "", true}, // max > allowed max 99 for segment specs
}
*/

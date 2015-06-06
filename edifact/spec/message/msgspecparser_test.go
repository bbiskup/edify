package message

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"sort"
	"testing"
	"time"
)

// Example of small file with 3 levels of nesting: TPFREP; multiple nesting levels end simulateously
// Bigger file, 3 levels, up/down: ORDRSP
// Most groups (258): GOVCBR; only msg type with > 99 groups

func TestParseINVOICFile(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	spec, err := parser.ParseSpecFile("../../../testdata/INVOIC_D.14B")
	assert.Nil(t, err)
	require.NotNil(t, spec)

	assert.Equal(t, "INVOIC", spec.Id)
	assert.Equal(t, "D", spec.Version)
	assert.Equal(t, "14B", spec.Release)
	assert.Equal(t, "UN", spec.ContrAgency)
	assert.Equal(t, "16", spec.Revision)
	assert.Equal(t, time.Date(2014, time.November, 17, 0, 0, 0, 0, time.UTC), spec.Date)

	fmt.Printf("Top level parts: " + spec.Dump())
	assert.Equal(t, 32, spec.Count())
	assert.Equal(t, "dummy_segment_spec-UNH", spec.TopLevelPart(0).Name())
	assert.Equal(t, "dummy_segment_spec-GIR", spec.TopLevelPart(10).Name())
	part11 := spec.TopLevelPart(11)
	assert.Equal(t, "Group_1", part11.Name())
	assert.Equal(t, 99999, part11.MaxCount())

	group_1, ok := part11.(*MsgSpecSegGrpPart)
	assert.True(t, ok)
	assert.Equal(t, group_1.Count(), 9)

	group_1_segm_0 := group_1.Children()[0].(*MsgSpecSegPart)
	assert.True(t, ok)
	assert.Equal(t, group_1_segm_0.SegSpec.Id, "RFF")
}

func TestParseAUTHORFile(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	spec, err := parser.ParseSpecFile("../../../testdata/AUTHOR_D.14B")
	assert.Nil(t, err)
	require.NotNil(t, spec)
	// fmt.Printf("Message spec: %s", spec)

	assert.Equal(t, "AUTHOR", spec.Id)
	assert.Equal(t, "D", spec.Version)
	assert.Equal(t, "14B", spec.Release)
	assert.Equal(t, "UN", spec.ContrAgency)
	assert.Equal(t, "3", spec.Revision)
	assert.Equal(t, time.Date(2014, time.November, 17, 0, 0, 0, 0, time.UTC), spec.Date)
}

func TestParseNonExistentFile(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	spec, err := parser.ParseSpecFile("../../../testdata/NON_EXISTENT")
	assert.NotNil(t, err)
	assert.Nil(t, spec)
}

func TestParseDir(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})

	specs, err := parser.ParseSpecDir("../../../testdata/message_specs", "14B")
	assert.Nil(t, err)
	require.NotNil(t, specs)
	assert.Equal(t, 2, len(specs))
	// fmt.Printf("Message specs: %s", specs)

	// cast necessary so sort.Interface methods will be recognized
	// on []*MsgSpec
	mSpecs := specs

	// ioutil.ReadDir sorts entries alphabetically
	balanc, ok := mSpecs["BALANC"]
	assert.True(t, ok)
	assert.Equal(t, "BALANC", balanc.Id)

	jobcon, ok := mSpecs["JOBCON"]
	assert.True(t, ok)
	assert.Equal(t, "JOBCON", jobcon.Id)
}

var segGrpStartSpec = []struct {
	line         string
	shouldMatch  bool
	recordNum    int
	groupNum     int
	isMandatory  bool
	maxCount     int
	nestingLevel int
}{
	{
		"00170       ---- Segment group 4  ------------------ C   99---------------+",
		true, 170, 4, false, 99, 1,
	},
	{
		"00130       ---- Segment group 3  ------------------ C   999-------------+|",
		true, 130, 3, false, 999, 2,
	},
	{
		"00210       ---- Segment group 6  ------------------ C   99-------------+||",
		true, 210, 6, false, 99, 3,
	},
	{
		"00050       ---- Segment group 1  ------------------ M   1----------------+",
		true, 50, 1, true, 1, 1,
	},
	{
		"00050       ---- Segment group 100------------------ M   1----------------+",
		true, 50, 100, true, 1, 1,
	},
	{
		// not a group start line
		"00090   RFF Reference                                C   9----------------+",
		false, 0, 0, false, 0, 1,
	},
	{
		// not a group start line
		"00110   EQD Equipment details                        M   1                |",
		false, 0, 0, false, 0, 1,
	},
}

func TestParseSegGrpStart(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	for _, spec := range segGrpStartSpec {
		res, err := parser.parseSegGrpStart(spec.line)
		require.Nil(t, err)

		if spec.shouldMatch {
			require.NotNil(t, res)
		} else {
			assert.Nil(t, res)
			continue
		}

		assert.Equal(t, spec.recordNum, res.RecordNum)
		assert.Equal(t, spec.groupNum, res.GroupNum)
		assert.Equal(t, spec.isMandatory, res.IsMandatory)
		assert.Equal(t, spec.maxCount, res.MaxCount)
		assert.Equal(t, spec.nestingLevel, res.NestingLevel)
	}
}

var segmentEntryStartSpec = []struct {
	line         string
	shouldMatch  bool
	recordNum    int
	segmentId    string
	segmentName  string
	isMandatory  bool
	maxCount     int
	nestingLevel int
}{
	{
		"00020   BGM Beginning of message                     M   1     ",
		true, 20, "BGM", "Beginning of message", true, 1, 0,
	},
	{
		"00060   TDT Transport information                    M   1                |",
		true, 60, "TDT", "Transport information", true, 1, 1,
	},
	{
		"00140   DTM Date/time/period                         M   1               ||",
		true, 140, "DTM", "Date/time/period", true, 1, 2,
	},
	{
		// not a segment entry
		"00210       ---- Segment group 6  ------------------ C   99-------------+||",
		false, 0, "", "", false, 0, 0,
	},
}

func TestParseSegEntry(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	for _, spec := range segmentEntryStartSpec {
		res, err := parser.parseSegEntry(spec.line)
		require.Nil(t, err)

		if spec.shouldMatch {
			require.NotNil(t, res)
		} else {
			assert.Nil(t, res)
			continue
		}
		assert.Equal(t, spec.recordNum, res.RecordNum)
		assert.Equal(t, spec.segmentId, res.SegId)
		assert.Equal(t, spec.segmentName, res.SegName)
		assert.Equal(t, spec.isMandatory, res.IsMandatory)
		assert.Equal(t, spec.maxCount, res.MaxCount)
		assert.Equal(t, spec.nestingLevel, res.NestingLevel)
	}
}

func TestParseHeaderSection(t *testing.T) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
	assert.Equal(t, true, parser.matchHeaderOrEmptyInGroupSection("            HEADER SECTION"))
	assert.Equal(t, true, parser.matchHeaderOrEmptyInGroupSection("            HEADER SECTION  "))
	assert.Equal(t, false, parser.matchHeaderOrEmptyInGroupSection("           HEADER SECTION"))
}

func BenchmarkParseDir(b *testing.B) {
	parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})

	for i := 0; i < b.N; i++ {
		//specs, err := parser.ParseSpecDir("../../../testdata/message_specs", "14B")
		specs, err := parser.ParseSpecDir("../../../testdata/d14b/edmd", "14B")
		assert.Nil(b, err)
		require.NotNil(b, specs)
		//assert.Equal(b, 2, len(specs))
		// fmt.Printf("Message specs: %s", specs)
	}
}

func BenchmarkParseFileTPFREP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
		spec, err := parser.ParseSpecFile("../../../testdata/d14b/edmd/TPFREP_D.14B")
		assert.Nil(b, err)
		require.NotNil(b, spec)
		require.NotNil(b, spec)
	}
}

func BenchmarkParseFileINVOIC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
		spec, err := parser.ParseSpecFile("../../../testdata/d14b/edmd/INVOIC_D.14B")
		assert.Nil(b, err)
		require.NotNil(b, spec)
		require.NotNil(b, spec)
	}
}

func BenchmarkParseFileGOVCBR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewMsgSpecParser(&MockSegSpecProviderImpl{})
		spec, err := parser.ParseSpecFile("../../../testdata/d14b/edmd/GOVCBR_D.14B")
		assert.Nil(b, err)
		require.NotNil(b, spec)
		require.NotNil(b, spec)
	}
}

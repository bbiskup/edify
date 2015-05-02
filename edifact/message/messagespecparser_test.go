package message

import (
	"fmt"
	"github.com/bbiskup/edify/edifact/segment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Example of small file with 3 levels of nesting: TPFREP; multiple nesting levels end simulateously
// Bigger file, 3 levels, up/down: ORDRSP
// Most groups (258): GOVCBR; only msg type with > 99 groups

func TestParseINVOICFile(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	spec, err := parser.ParseSpecFile("../../testdata/INVOIC_D.14B")
	assert.Nil(t, err)
	assert.NotNil(t, spec)
	fmt.Printf("Message spec: %s", spec)

	assert.Equal(t, "INVOIC", spec.Id)
	assert.Equal(t, "D", spec.Version)
	assert.Equal(t, "14B", spec.Release)
	assert.Equal(t, "UN", spec.ContrAgency)
	assert.Equal(t, "16", spec.Revision)
	assert.Equal(t, time.Date(2014, time.November, 17, 0, 0, 0, 0, time.UTC), spec.Date)
}

func TestParseAUTHORFile(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	spec, err := parser.ParseSpecFile("../../testdata/AUTHOR_D.14B")
	assert.Nil(t, err)
	assert.NotNil(t, spec)
	fmt.Printf("Message spec: %s", spec)

	assert.Equal(t, "AUTHOR", spec.Id)
	assert.Equal(t, "D", spec.Version)
	assert.Equal(t, "14B", spec.Release)
	assert.Equal(t, "UN", spec.ContrAgency)
	assert.Equal(t, "3", spec.Revision)
	assert.Equal(t, time.Date(2014, time.November, 17, 0, 0, 0, 0, time.UTC), spec.Date)
}

func TestParseNonExistentFile(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	spec, err := parser.ParseSpecFile("../../testdata/NON_EXISTENT")
	assert.NotNil(t, err)
	assert.Nil(t, spec)
}

func TestParseDir(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	specs, err := parser.ParseSpecDir("../../testdata/message_specs", "14B")
	assert.Nil(t, err)
	assert.NotNil(t, specs)
	fmt.Printf("Message specs: %s", specs)

	// ioutil.ReadDir sorts entries alphabetically
	assert.Equal(t, "BALANC", specs[0].Id)
	assert.Equal(t, "JOBCON", specs[1].Id)
}

var segmentGroupStartSpec = []struct {
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

func TestParseSegmentGroupStart(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	for _, spec := range segmentGroupStartSpec {
		res, err := parser.parseSegmentGroupStart(spec.line)
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
		true, 60, "TDT", "Transport information", true, 1, 1},
}

func TestParseSegmentEntry(t *testing.T) {
	parser := NewMessageSpecParser(segment.SegmentSpecMap{})
	for _, spec := range segmentEntryStartSpec {
		res, err := parser.parseSegmentEntry(spec.line)
		require.Nil(t, err)

		if spec.shouldMatch {
			require.NotNil(t, res)
		} else {
			assert.Nil(t, res)
			continue
		}
		assert.Equal(t, spec.recordNum, res.RecordNum)
		assert.Equal(t, spec.segmentId, res.SegmentId)
		assert.Equal(t, spec.segmentName, res.SegmentName)
		assert.Equal(t, spec.isMandatory, res.IsMandatory)
		assert.Equal(t, spec.maxCount, res.MaxCount)
		assert.Equal(t, spec.nestingLevel, res.NestingLevel)
	}
}

package message

type SegEntry struct {
	RecordNum   int
	SegId       string
	SegName     string
	IsMandatory bool
	MaxCount    int

	// Nesting level _after_ this segment entry
	// A segment entry might close multiple groups simultaneously.
	NestingLevel int
}

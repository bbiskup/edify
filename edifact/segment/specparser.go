package segment

// Parses segment specifications file (e.g. EDSD.14B)
type SegmentSpecParser struct {
}

type SegmentSpecMap map[string]SegmentSpec

func (p *SegmentSpecParser) ParseSpecFile(fileName string) (specs SegmentSpecMap, err error) {
	panic("not implemented")
}

func NewSegmentSpecParser() *SegmentSpecParser {
	return &SegmentSpecParser{}
}

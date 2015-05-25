package msg

// A group of segments. Groups are not identified during message parsing,
// but during validation, when message specification is available
type SegmentGroup struct {
	id    string
	Parts []SegmentOrGroup
}

// From interface SegmentOrGroup
func (g *SegmentGroup) Id() string {
	return g.id
}

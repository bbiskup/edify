package validation

type SegSeqState int

const (
	seqStateInitial SegSeqState = iota
	seqStateGroupStart
	seqStateSeg
	seqStateSearching
)

func (st SegSeqState) String() string {
	switch st {
	case seqStateInitial:
		return "initial"
	case seqStateGroupStart:
		return "group_start"
	case seqStateSeg:
		return "seg"
	case seqStateSearching:
		return "searching"
	default:
		panic("Unknown state")
	}
}

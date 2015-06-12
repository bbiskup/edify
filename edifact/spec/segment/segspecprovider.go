package segment

// TODO import spec from r1241.txt
func IsUnValidatedSegment(segID string) bool {
	return segID == "UNH" || segID == "UNT" || segID == "UNS" || segID == "UGH" || segID == "UGT"
}

// Provides segment spec by Id
type SegSpecProvider interface {
	Get(id string) *SegSpec
	Len() int
}

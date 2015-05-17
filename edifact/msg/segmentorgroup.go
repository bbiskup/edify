package msg

// For storing segments and groups in the same array
type SegmentOrGroup interface {
	Id() string
}

package msg

// For storing segments and groups in the same array
type SegOrGroup interface {
	MsgDumper
	Id() string
}

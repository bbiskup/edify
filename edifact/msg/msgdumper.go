package msg

// Human-readable output of nested message structure
type MsgDumper interface {
	Dump(indent int) string
}

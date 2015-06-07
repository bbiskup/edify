package rawmsg

import (
	"strings"
)

// Human-readable output of nested message structure
type MsgDumper interface {
	Dump(indent int) string
}

func getIndentStr(indent int) string {
	return strings.Repeat("  ", indent)
}

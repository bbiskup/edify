package msg

// Common interface of constituents of validated messages
type NestedMsgPart interface {
	Id() string
}

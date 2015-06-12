package msg

type RepMsgPart interface {
	MsgDumper
	Count() int
	Id() string
}

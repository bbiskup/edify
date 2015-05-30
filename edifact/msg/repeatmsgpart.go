package msg

type RepeatMsgPart interface {
	MsgDumper
	Count() int
	Id() string
}

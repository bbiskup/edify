package msg

type Interchange struct {
	RawMessages []*RawMessage
}

func (i *Interchange) String() string {
	result := ""
	for _, m := range i.RawMessages {
		result += "\n" + m.String() + "\n"
	}
	return result
}

func (i *Interchange) AddMessage(message *RawMessage) {
	i.RawMessages = append(i.RawMessages, message)
}

func NewInterchange() *Interchange {
	return &Interchange{}
}

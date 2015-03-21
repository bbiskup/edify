package edifact

type Interchange struct {
	Messages []*Message
}

func (i *Interchange) String() string {
	result := ""
	for _, m := range i.Messages {
		result += "\n" + m.String() + "\n"
	}
	return result
}

func (i *Interchange) AddMessage(message *Message) {
	i.Messages = append(i.Messages, message)
}

func NewInterchange() *Interchange {
	return &Interchange{}
}

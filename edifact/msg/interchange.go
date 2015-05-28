package msg

import (
	"bytes"
	"fmt"
)

type Interchange struct {
	RawMessages []*RawMessage
}

func (i *Interchange) String() string {
	var buf bytes.Buffer
	for _, m := range i.RawMessages {
		buf.WriteString(fmt.Sprintf("\n%s\n", m.String()))
	}
	return buf.String()
}

func (i *Interchange) AddMessage(message *RawMessage) {
	i.RawMessages = append(i.RawMessages, message)
}

func NewInterchange() *Interchange {
	return &Interchange{}
}

package msg

import (
	"bytes"
	"fmt"
)

type Interchange struct {
	RawMsgs []*RawMsg
}

func (i *Interchange) String() string {
	var buf bytes.Buffer
	for _, m := range i.RawMsgs {
		buf.WriteString(fmt.Sprintf("\n%s\n", m.String()))
	}
	return buf.String()
}

func (i *Interchange) AddMessage(message *RawMsg) {
	i.RawMsgs = append(i.RawMsgs, message)
}

func NewInterchange() *Interchange {
	return &Interchange{}
}

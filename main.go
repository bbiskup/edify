package main

import (
	edi "edifact_experiments/edifact"
	"fmt"
)

func main() {
	e := edi.NewElement("name1", "value1")
	s := edi.NewSegment("segname1")
	s.AddElement(e)

	m := edi.NewMessage("mssagename1")
	m.AddSegment(s)

	i := edi.NewInterchange()
	i.AddMessage(m)

	fmt.Println(i)
}

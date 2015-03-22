package main

import (
	//edi "edifact_experiments/edifact"
	sp "edifact_experiments/edifact/specparser"
	"fmt"
)

func main() {
	/*
		e := edi.NewElement("name1", "value1")
		s := edi.NewSegment("segname1")
		s.AddElement(e)

		m := edi.NewMessage("messagename1")
		m.AddSegment(s)

		i := edi.NewInterchange()
		i.AddMessage(m)

		fmt.Println(i)
	*/

	p := sp.NewDataElementSpecParser()
	specs, err := p.ParseSpecFile("testdata/EDED.14B_short")
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		return
	}
	fmt.Printf("Specs: %s\n", specs)
}

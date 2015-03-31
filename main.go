package main

import (
	"fmt"
	sp "github.com/bbiskup/edifice/edifact/dataelement"
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
	specs, err := p.ParseSpecFile("testdata/EDED.14B")
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		return
	}
	fmt.Printf("Specs:\n")
	for _, spec := range specs {
		fmt.Printf("\t%s\n", spec)
	}
}

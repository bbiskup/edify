package commands

import (
	"errors"
	"fmt"
	sp "github.com/bbiskup/edify/edifact/dataelement"
	"log"
	"os"
	"strings"
)

func Parse(fileName string) error {
	if len(fileName) == 0 {
		return errors.New("No filename given")
	}
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if strings.Index(fileName, "EDED") == 0 {
		return ParseSimpleDataElements(fileName)
	}
	return errors.New(fmt.Sprintf("Unrecognized file: %s", fileName))
}

func ParseSimpleDataElements(fileName string) error {
	log.Printf("ParseSimpleDataElements %s\n", fileName)
	p := sp.NewSimpleDataElementSpecParser()
	specs, err := p.ParseSpecFile(fileName)
	if err != nil {
		return err
	}
	fmt.Printf("Specs:\n")
	for _, spec := range specs {
		fmt.Printf("\t%s\n", spec)
	}
	fmt.Println("")
	return nil
}

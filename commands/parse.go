package commands

import (
	"errors"
	"fmt"
	"github.com/bbiskup/edify/edifact/spec/codes"
	// "github.com/bbiskup/edify/edifact/dataelement"
	"log"
	"os"
	"strings"
)

func Parse(fileNames []string) error {
	if len(fileNames) == 0 {
		return errors.New("Nothing to parse")
	}
	for _, fileName := range fileNames {
		err := ParseFile(fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseFile(fileName string) error {
	if len(fileName) == 0 {
		return errors.New("No filename given")
	}
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	pathParts := strings.Split(fileName, string(os.PathSeparator))
	filePart := pathParts[len(pathParts)-1]

	/*if strings.HasPrefix(filePart, "EDED") {
		return ParseSimpleDataElems(fileName)
	}

	if strings.HasPrefix(filePart, "EDCD") {
		return ParseCompositeDataElems(fileName)
	}*/

	if strings.HasPrefix(filePart, "UNCL") {
		_, err := ParseCodeList(fileName)
		return err
	}

	return errors.New(fmt.Sprintf("Unrecognized file: %s", fileName))
}

// func ParseSimpleDataElems(fileName string) error {
// 	log.Printf("ParseSimpleDataElems %s\n", fileName)
// 	p := dataelement.NewSimpleDataElemSpecParser()
// 	specs, err := p.ParseSpecFile(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Found %d specs", len(specs))
// 	fmt.Println("")
// 	return nil
// }

// func ParseCompositeDataElems(fileName string) error {
// 	log.Printf("ParseCompositeDataElems %s\n", fileName)
// 	p := dataelement.NewCompositeDataElemSpecParser()
// 	specs, err := p.ParseSpecFile(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Found %d specs", len(specs))
// 	fmt.Println("")
// 	return nil
// }

func ParseCodeList(fileName string) (codes.CodesSpecMap, error) {
	log.Printf("ParseCodeList %s\n", fileName)
	p := codes.NewCodesSpecParser()
	specs, err := p.ParseSpecFile(fileName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Found %d specs", len(specs))
	/*fmt.Printf("Specs:\n")
	for _, spec := range specs {
		fmt.Printf("\t%s\n", spec)
	}*/
	fmt.Println("")
	return specs, nil
}

package commands

import (
	"errors"
	"fmt"
	csp "github.com/bbiskup/edify/edifact/spec/codes"
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

	if strings.HasPrefix(filePart, "UNCL") {
		_, err := ParseCodeList(fileName)
		return err
	}

	return fmt.Errorf("Unrecognized file: %s", fileName)
}

func ParseCodeList(fileName string) (csp.CodesSpecMap, error) {
	log.Printf("ParseCodeList %s\n", fileName)
	p := csp.NewCodesSpecParser()
	specs, err := p.ParseSpecFile(fileName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Found %d specs", len(specs))
	fmt.Println("")
	return specs, nil
}

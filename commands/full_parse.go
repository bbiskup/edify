package commands

import (
	"errors"
	"github.com/bbiskup/edify/edifact/spec/specparser"
	"github.com/bbiskup/edify/edifact/validation"
	"log"
)

func FullParse(version string, specDirName string) error {
	if version == "" {
		return errors.New("No version given")
	}

	if specDirName == "" {
		return errors.New("No spec dir given given")
	}

	parser, err := specparser.NewFullSpecParser(version, specDirName)
	if err != nil {
		return err
	}
	segSpecs, err := parser.ParseSegSpecsWithPrerequisites()
	if err != nil {
		return err
	}
	msgSpecs, err := parser.ParseMsgSpecs(segSpecs)
	if err != nil {
		return err
	}
	validator := validation.NewMsgValidator(msgSpecs, segSpecs)

	log.Printf("Validator has %d message specs and %d segment specs",
		validator.MsgSpecCount(), validator.SegSpecCount())
	return nil
}

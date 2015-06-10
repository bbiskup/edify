package commands

import (
	"github.com/bbiskup/edify/edifact/validation"
	"log"
)

func FullParse(version string, specDirName string) error {
	validator, err := validation.GetMsgValidator(version, specDirName)
	if err != nil {
		return err
	}

	log.Printf("Validator has %d message specs and %d segment specs",
		validator.MsgSpecCount(), validator.SegSpecCount())
	return nil
}

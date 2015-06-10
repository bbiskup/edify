package commands

import (
	"errors"
	"github.com/bbiskup/edify/edifact/query"
	"github.com/bbiskup/edify/edifact/rawmsg"
	"github.com/bbiskup/edify/edifact/validation"
	"log"
)

func Query(version string, specDirName string, msgFileName string, queryStr string) error {
	if version == "" {
		return errors.New("No version specified")
	}

	if specDirName == "" {
		return errors.New("No specification directory specified")
	}

	if msgFileName == "" {
		return errors.New("No message file specified")
	}

	validator, err := validation.GetMsgValidator(version, specDirName)
	if err != nil {
		return err
	}

	rawMsgParser := rawmsg.NewParser()
	rawMsg, err := rawMsgParser.ParseRawMsgFile(msgFileName)
	if err != nil {
		return err
	}

	nestedMsg, err := validator.Validate(rawMsg)
	if err != nil {
		return err
	}
	log.Printf("Nested msg: %s", nestedMsg.Dump())

	if queryStr != "" {
		navigator := query.NewNavigator()
		queryResult, err := navigator.Navigate(queryStr, nestedMsg)
		if err != nil {
			return err
		}
		log.Printf("Query result: %s", queryResult)
	} else {
		log.Printf("No query string specified")
	}
	return nil
}

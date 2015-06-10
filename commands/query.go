package commands

import (
	"errors"
	"github.com/bbiskup/edify/edifact/query"
	"github.com/bbiskup/edify/edifact/rawmsg"
	"github.com/bbiskup/edify/edifact/validation"
	"log"
)

func checkParams(version string, specDirName string,
	msgFileName string, dumpMessage bool,
	queryStr string) error {

	if version == "" {
		return errors.New("No version specified")
	}

	if specDirName == "" {
		return errors.New("No specification directory specified")
	}

	if queryStr != "" && msgFileName == "" {
		return errors.New("Query not possible; no message file specified")
	}

	if dumpMessage && msgFileName == "" {
		return errors.New("Message dump not possible; no message file specified")
	}
	return nil
}

func Query(
	version string, specDirName string,
	msgFileName string, dumpMessage bool,
	queryStr string) error {

	if err := checkParams(version, specDirName, msgFileName, dumpMessage, queryStr); err != nil {
		return err
	}

	validator, err := validation.GetMsgValidator(version, specDirName)
	if err != nil {
		return err
	}

	var rawMsg *rawmsg.RawMsg
	if msgFileName != "" {
		rawMsgParser := rawmsg.NewParser()
		rawMsg, err = rawMsgParser.ParseRawMsgFile(msgFileName)
		if err != nil {
			return err
		}
	} else {
		log.Printf("No message file specified")
		return nil
	}

	nestedMsg, err := validator.Validate(rawMsg)
	if err != nil {
		return err
	}

	if dumpMessage {
		log.Printf("Nested msg: %s", nestedMsg.Dump())
	}

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

package commands

import (
	"github.com/bbiskup/edify/edifact/query"
	"github.com/bbiskup/edify/edifact/rawmsg"
	"github.com/bbiskup/edify/edifact/validation"
	"log"
)

func Query(version string, specDirName string, msgFileName string, queryStr string) error {
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

	navigator := query.NewNavigator()
	queryResult, err := navigator.Navigate(queryStr, nestedMsg)
	if err != nil {
		return err
	}
	log.Printf("Query result: %s", queryResult)

	return nil
}

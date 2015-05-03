package commands

import (
	"errors"
	"github.com/bbiskup/edify/edifact/spec/specparser"
)

func FullParse(specDirNames []string) error {
	if len(specDirNames) == 0 {
		return errors.New("No spec directory names given")
	}

	for _, specDirName := range specDirNames {
		parser, err := specparser.NewFullSpecParser("14B", specDirName)
		if err != nil {
			return err
		}

		err = parser.Parse()
		if err != nil {
			return err
		}
	}
	return nil
}

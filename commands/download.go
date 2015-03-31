package commands

import (
	"errors"
	"log"
)

const (
	downloadDir = "TODO"
)

func Download(url string) error {
	log.Printf("Download %s", url)
	if len(url) == 0 {
		return errors.New("No URL specified")
	}
	return nil
}

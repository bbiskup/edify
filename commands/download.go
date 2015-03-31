package commands

import (
	"log"
)

const (
	downloadDir = "TODO"
)

func Download(url string) error {
	log.Printf("Download %s", url)
	return nil
}

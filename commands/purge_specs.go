package commands

import (
	"errors"
	"log"
	"os"
)

// Purge extracted EDIFACT spec content
func PurgeSpecs(version string) error {
	if len(version) == 0 {
		return errors.New("No version specified")
	}
	versionSubDir := downloadDir + string(os.PathSeparator) + versionDir(version)
	s, err := os.Stat(versionSubDir)
	if !os.IsNotExist(err) && s.IsDir() {
		log.Printf("Purging version %s (%s)", version, versionSubDir)
		err = os.RemoveAll(versionSubDir)
		return err
	} else {
		log.Printf("%s does not exist; not need to purge", versionSubDir)
	}
	return nil
}

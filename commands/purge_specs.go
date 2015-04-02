package commands

import (
	"log"
	"os"
)

// Purge extracted EDIFACT spec content
func PurgeSpecs(version string) error {
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

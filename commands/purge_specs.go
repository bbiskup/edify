package commands

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const pathSep = string(os.PathSeparator)

// Delete all .zip files in download dir
func purgeDownloads() error {
	entries, err := ioutil.ReadDir(downloadDir)
	if err != nil {
		return err
	}
	log.Printf("Purging %d downloaded archives in %s", len(entries), downloadDir)

	for _, entry := range entries {
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".zip") {
			log.Printf("Skipping %s", fileName)
			continue
		}
		log.Printf("Deleting %s", fileName)
		err := os.Remove(downloadDir + pathSep + fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Purge extracted EDIFACT spec content, and possibly
// download archives if purgeAll == true
func PurgeSpecs(version string, purgeAll bool) error {
	if purgeAll {
		err := purgeDownloads()
		if err != nil {
			return err
		}
	}

	if len(version) == 0 {
		return errors.New("No version specified")
	}
	versionSubDir := downloadDir + pathSep + versionDir(version)
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

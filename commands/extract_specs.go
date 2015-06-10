package commands

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func extractSpecsFirstLevel(archivePath string, targetDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create target dir

	var count int32
	for _, f := range reader.File {
		log.Printf("Extracting %s", f.Name)
		contents, err := f.Open()
		if err != nil {
			return err
		}
		defer contents.Close()

		targetFile, err := os.Create(targetDir + string(os.PathSeparator) + f.Name)
		if err != nil {
			return err
		}

		_, err = io.Copy(targetFile, contents)
		if err != nil {
			return err
		}
		count++
	}
	log.Printf("Extracted %d files on 1st level", count)
	return nil
}

// Extract zips inside top-level spec zip to same dir
func extractInnerZIP(targetDir string, archiveFile string) error {
	archivePath := targetDir + string(os.PathSeparator) + archiveFile

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}

	archiveFileTokens := strings.Split(archiveFile, ".")
	archiveFilePrefix := archiveFileTokens[0]

	var count int32
	for _, f := range reader.File {
		log.Printf("Extracting %s (2nd level)", f.Name)
		contents, err := f.Open()
		if err != nil {
			return err
		}
		defer contents.Close()

		targetSubDir := strings.Join([]string{
			targetDir, archiveFilePrefix,
		}, string(os.PathSeparator))

		err = os.MkdirAll(targetSubDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}

		targetFile, err := os.Create(targetSubDir + string(os.PathSeparator) + f.Name)
		if err != nil {
			return err
		}

		_, err = io.Copy(targetFile, contents)
		if err != nil {
			return err
		}

		count++
	}
	log.Printf("Extracted %d files on 1st level", count)

	// Inner zip is no longer needed
	err = os.Remove(archivePath)
	if err != nil {
		return err
	}

	return nil
}

// Extract zips inside top-level spec zip
func extractSpecsSecondLevel(targetDir string) error {
	dirContents, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return err
	}
	for _, entry := range dirContents {
		name := entry.Name()
		if !strings.HasSuffix(name, ".zip") {
			continue
		}
		log.Printf("Extracting %s (%.2f MB)", name, float32(entry.Size())/1e6)

		err = extractInnerZIP(targetDir, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Extract EDIFACT specs from file provided by UNECE
func ExtractSpecs(version string) error {
	if len(version) == 0 {
		return errors.New("No version specified")
	}
	version = strings.ToLower(version)

	archivePath := downloadPath(version)
	targetDir := downloadDir + string(os.PathSeparator) + versionDir(version)

	log.Printf("Extracting archive %s --> %s", archivePath, targetDir)

	err := os.MkdirAll(targetDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	err = extractSpecsFirstLevel(archivePath, targetDir)
	if err != nil {
		return err
	}

	return extractSpecsSecondLevel(targetDir)
}

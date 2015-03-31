package commands

import (
	"archive/zip"
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

	var count int32
	for _, f := range reader.File {
		log.Printf("Extracting %s (2nd level)", f.Name)
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
	archivePath := downloadPath(version)

	targetDir := downloadDir + string(os.PathSeparator) + version
	err := os.MkdirAll(targetDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	err = extractSpecsFirstLevel(archivePath, targetDir)
	if err != nil {
		return err
	}

	err = extractSpecsSecondLevel(targetDir)

	return nil
}

package commands

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

func extractSpecsFirstLevel(archivePath string, targetDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create target dir

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
	}
	return nil
}

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
	return nil
}

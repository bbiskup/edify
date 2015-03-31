package commands

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

func ExtractSpecs(version string) error {
	archivePath := downloadPath(version)

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create target dir
	targetDirName := downloadDir + string(os.PathSeparator) + version
	err = os.Mkdir(targetDirName, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		log.Printf("Extracting %s", f.Name)
		contents, err := f.Open()
		if err != nil {
			return err
		}
		defer contents.Close()

		targetFile, err := os.Create(targetDirName + string(os.PathSeparator) + f.Name)
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

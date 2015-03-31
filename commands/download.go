package commands

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	downloadDir = ".edify/downloads"
	specFile    = "edififact_spec"
)

// Retrieve EDIFACT specification files for validation
func DownloadSpecs(urlStr string) error {
	log.Printf("Download %s", urlStr)
	if len(urlStr) == 0 {
		return errors.New("No URL specified")
	}

	// Validate URL
	_, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	targetPath := downloadDir + string(os.PathSeparator) + specFile

	// Ensure the download directory exists
	err = os.MkdirAll(downloadDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	response, err := check.Get(urlStr)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	log.Printf("Download status: %s", response.Status)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Download failed with status %s", response.Status))
	}

	size, err := io.Copy(targetFile, response.Body)
	if err != nil {
		return err
	}

	log.Printf("Download of '%s' (%f bytes) complete", urlStr, size/1e6)

	return nil
}

package commands

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	urlRoot     = "http://www.unece.org/fileadmin/DAM/trade/untdid/"
	downloadDir = ".edify/downloads"
)

func prepareTargetPath(version string) (*os.File, error) {
	// Ensure the download directory exists
	err := os.MkdirAll(downloadDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}

	targetPath := downloadDir + string(os.PathSeparator) + version

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}

	return targetFile, nil
}

// Retrieve EDIFACT specification files for validation
func DownloadSpecs(version string) error {

	if len(version) == 0 {
		return errors.New("No version specified")
	}

	urlStr := strings.Join([]string{urlRoot, version, version + ".zip"}, "/")

	// Validate URL
	_, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	targetFile, err := prepareTargetPath(version)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	log.Printf("Download version %s (%s)", version, urlStr)

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

	log.Printf("Download of '%s' (%.2f MB) complete", urlStr, float64(size)/1e6)

	return nil
}

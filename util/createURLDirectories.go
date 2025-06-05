package util

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// parse url
// get host
// get path components excluding the file name
// split path and exclude last component if its a file
// create the directories
func CreateURLDirectories(rawURL string, baseDir string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("failed to parse URL %s: %v", rawURL, err)
		return "", err
	}

	host := parsedURL.Host
	if host == "" {

		pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(pathSegments) > 0 {
			host = pathSegments[0]
		} else {
			host = "unknown"
		}
	}

	dirPath := strings.Trim(parsedURL.Path, "/")
	if dirPath == "" {

		dir := filepath.Join(baseDir, host)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			err = fmt.Errorf("failed to create directory %s: %v", dir, err)
			return "", err
		}
		return dir, nil
	}

	pathSegments := strings.Split(dirPath, "/")
	dirSegments := pathSegments
	if len(pathSegments) > 0 {

		lastSegment := pathSegments[len(pathSegments)-1]
		if strings.Contains(lastSegment, ".") {
			dirSegments = pathSegments[:len(pathSegments)-1]
		}
	}

	dir := filepath.Join(baseDir, host, filepath.Join(dirSegments...))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		err = fmt.Errorf("failed to create directory %s: %v", dir, err)
		return "", err
	}

	return dir, nil
}

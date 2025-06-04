package util

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func CreateURLDirectories(rawURL string, baseDir string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("failed to parse URL %s: %v", rawURL, err)
		return "", err
	}

	// Get the host
	host := parsedURL.Host
	if host == "" {
		// If no host, use the first path segment or a default
		pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(pathSegments) > 0 {
			host = pathSegments[0]
		} else {
			host = "unknown"
		}
	}

	// Get the path components, excluding the file name
	dirPath := strings.Trim(parsedURL.Path, "/")
	if dirPath == "" {
		// Only host, no path
		dir := filepath.Join(baseDir, host)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			err = fmt.Errorf("failed to create directory %s: %v", dir, err)
			return "", err
		}
		return dir, nil
	}

	// Split path and exclude the last component if it’s a file
	pathSegments := strings.Split(dirPath, "/")
	dirSegments := pathSegments
	if len(pathSegments) > 0 {
		// Check if the last segment is a file (has an extension)
		lastSegment := pathSegments[len(pathSegments)-1]
		if strings.Contains(lastSegment, ".") {
			dirSegments = pathSegments[:len(pathSegments)-1]
		}
	}

	// Construct the directory path
	dir := filepath.Join(baseDir, host, filepath.Join(dirSegments...))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		err = fmt.Errorf("failed to create directory %s: %v", dir, err)
		return "", err
	}

	return dir, nil
}

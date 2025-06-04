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

	

	

	
}

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

	

	

	

	
}

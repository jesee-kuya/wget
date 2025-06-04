package parser

import (
	"fmt"
	"net/url"
	"os"
)

// RewriteLinksInHTML rewrites internal URLs in an HTML file to local relative paths.
// It takes the file path, base URL, and root directory where files are saved.
func RewriteLinksInHTML(filePath string, baseURL *url.URL, rootDir string) error {
	// Read the HTML file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open HTML file %s: %w", filePath, err)
	}
	defer file.Close()

	return nil
}

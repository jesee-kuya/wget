package parser

import (
	"fmt"
	"net/url"
	"os"

	"golang.org/x/net/html"
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

	// Parse the HTML
	doc, err := html.Parse(file)
	if err != nil {
		return fmt.Errorf("failed to parse HTML file %s: %w", filePath, err)
	}

	// Rewrite links
	if err := rewriteNode(doc, baseURL, rootDir, filePath); err != nil {
		return fmt.Errorf("failed to rewrite links in %s: %w", filePath, err)
	}

	return nil
}

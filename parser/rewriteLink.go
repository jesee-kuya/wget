package parser

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

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

	// Create a temporary file for the rewritten content
	tempFile, err := os.CreateTemp(filepath.Dir(filePath), "temp_rewrite_*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file for %s: %w", filePath, err)
	}
	defer tempFile.Close()

	// Write the modified HTML to the temp file
	if err := html.Render(tempFile, doc); err != nil {
		return fmt.Errorf("failed to render HTML for %s: %w", filePath, err)
	}

	// Replace the original file with the temp file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file for %s: %w", filePath, err)
	}
	if err := os.Rename(tempFile.Name(), filePath); err != nil {
		return fmt.Errorf("failed to replace original file %s: %w", filePath, err)
	}

	return nil
}

// rewriteNode traverses the HTML node tree and rewrites internal URLs to local paths.
func rewriteNode(n *html.Node, baseURL *url.URL, rootDir, filePath string) error {
	if n.Type == html.ElementNode {
		var attrKey string
		switch n.Data {
		case "a", "link":
			attrKey = "href"
		case "img", "script", "iframe":
			attrKey = "src"
		}

		if attrKey != "" {
			for i, attr := range n.Attr {
				if attr.Key == attrKey && attr.Val != "" {
					parsedURL, err := baseURL.Parse(strings.TrimSpace(attr.Val))
					if err != nil {
						continue // Skip invalid URLs
					}

					// Only rewrite internal URLs (same host)
					if parsedURL.Host == baseURL.Host {
						// Get the relative path from the root directory
						localPath, err := convertURLToLocalPath(parsedURL, rootDir)
						if err != nil {
							return err
						}


	return nil
}


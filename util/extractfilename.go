package util

import (
	"net/url"
	"path"
	"strings"
)

// ExtractFilenameFromURL returns the base filename from a URL.
// If no filename is found, it falls back to "index.html".
func ExtractFilenameFromURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Path == "" {
		return "index.html"
	}
	filename := path.Base(parsed.Path)
	if filename == "" || strings.HasSuffix(filename, "/") {
		return "index.html"
	}
	return filename
}

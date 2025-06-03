package parser

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

// ExtractLinks parses HTML from the reader and returns internal asset URLs.
// It extracts href/src from <a>, <link>, and <img> tags.
func ExtractLinks(baseURL *url.URL, r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []string
	visited := make(map[string]bool)

	return links, nil
}

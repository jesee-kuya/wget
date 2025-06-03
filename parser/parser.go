package parser

import (
	"io"
	"net/url"
	"strings"

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

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var key string
			switch n.Data {
			case "a", "link":
				key = "href"
			case "img":
				key = "src"
			}

			if key != "" {
				for _, attr := range n.Attr {
					if attr.Key == key {
						href := strings.TrimSpace(attr.Val)
						if href == "" {
							continue
						}
						u, err := baseURL.Parse(href)
						if err == nil && u.Host == baseURL.Host {
							cleanURL := u.String()
							if !visited[cleanURL] {
								visited[cleanURL] = true
								links = append(links, cleanURL)
							}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return links, nil
}

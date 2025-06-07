package parser

import (
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// regex to capture the contents of url(...)
var cssURLRe = regexp.MustCompile(`url\(\s*['"]?([^)'"]+)['"]?\s*\)`)

// ExtractLinks parses HTML from the reader and returns internal asset URLs.
// It extracts href/src from <a>, <link>, <img> tags and url(...) in CSS.
func ExtractLinks(baseURL *url.URL, r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []string
	visited := make(map[string]bool)

	addLink := func(raw string) {
		href := strings.TrimSpace(raw)
		if href == "" {
			return
		}
		u, err := baseURL.Parse(href)
		if err != nil || u.Host != baseURL.Host {
			return
		}
		clean := u.String()
		if !visited[clean] {
			visited[clean] = true
			links = append(links, clean)
		}
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 1) handle <a>, <link>, <img>
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
						addLink(attr.Val)
					}
				}
			}

			// 2) handle <style> tags
			if n.Data == "style" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				css := n.FirstChild.Data
				matches := cssURLRe.FindAllStringSubmatch(css, -1)
				for _, m := range matches {
					addLink(m[1])
				}
			}

			// 3) handle inline style attributes
			for _, attr := range n.Attr {
				if attr.Key == "style" {
					css := attr.Val
					matches := cssURLRe.FindAllStringSubmatch(css, -1)
					for _, m := range matches {
						addLink(m[1])
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

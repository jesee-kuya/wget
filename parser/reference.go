package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// RewriteLinks parses the HTML in inBuf, rewrites all internal <a>, <img>, and <link> URLs
// to point at their local paths (rootDir/<domain>), and returns the modified HTML.
func RewriteLinks(inBuf []byte, pageURL *url.URL, rootDir string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(inBuf))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

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
				for i, attr := range n.Attr {
					if attr.Key == key {
						orig := strings.TrimSpace(attr.Val)
						if orig == "" {
							continue
						}
						u, err := pageURL.Parse(orig)
						if err != nil || u.Host != pageURL.Host {
							continue
						}
						localPath := filepath.ToSlash(u.Path)
						n.Attr[i].Val = localPath
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	var out bytes.Buffer
	if err := html.Render(&out, doc); err != nil {
		return nil, fmt.Errorf("rendering HTML: %w", err)
	}
	return out.Bytes(), nil
}

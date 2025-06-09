package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// RewriteLinks parses the HTML in inBuf, and for each internal link (a[href], img[src], link[href], style attributes),
// rewrites it to be relative to the HTML file’s location.
func RewriteLinks(inBuf []byte, pageURL *url.URL, domainDir, htmlDir string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(inBuf))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	// Rewrite CSS url(...) to relative path
	replaceCSSURLs := func(css string) string {
		return cssURLRe.ReplaceAllStringFunc(css, func(match string) string {
			matches := cssURLRe.FindStringSubmatch(match)
			if len(matches) < 2 {
				return match
			}
			orig := strings.TrimSpace(matches[1])
			u, err := pageURL.Parse(orig)
			if err != nil || u.Host != pageURL.Host {
				return match
			}
			localAbs := filepath.Join(domainDir, filepath.FromSlash(u.Path))
			rel, err := filepath.Rel(htmlDir, localAbs)
			if err != nil {
				rel = filepath.ToSlash(u.Path)
			} else {
				rel = filepath.ToSlash(rel)
			}
			return fmt.Sprintf("url('%s')", rel)
		})
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Handle standard tags: a[href], link[href], img[src]
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
						localAbs := filepath.Join(domainDir, filepath.FromSlash(u.Path))
						rel, err := filepath.Rel(htmlDir, localAbs)
						if err != nil {
							rel = filepath.ToSlash(u.Path)
						} else {
							rel = filepath.ToSlash(rel)
						}
						n.Attr[i].Val = rel
					}
				}
			}

			// Handle inline style="..."
			for i, attr := range n.Attr {
				if attr.Key == "style" && strings.Contains(attr.Val, "url(") {
					n.Attr[i].Val = replaceCSSURLs(attr.Val)
				}
			}

			// Handle <style> blocks
			if n.Data == "style" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				n.FirstChild.Data = replaceCSSURLs(n.FirstChild.Data)
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

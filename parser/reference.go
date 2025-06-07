package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// RewriteLinks parses the HTML in inBuf, and for each internal link (a[href], img[src], link[href]),
// rewrites it to be relative to the HTML file’s location.
// Parameters:
//   - inBuf: original HTML bytes
//   - pageURL: the URL of the HTML page
//   - domainDir: local root directory where the domain’s files are stored
//   - htmlDir: the local directory of this HTML file (i.e. filepath.Dir(outputPath))
func RewriteLinks(inBuf []byte, pageURL *url.URL, domainDir, htmlDir string) ([]byte, error) {
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
						// Local asset absolute path on disk:
						localAbs := filepath.Join(domainDir, filepath.FromSlash(u.Path))
						// Compute relative path from this HTML file’s directory:
						rel, err := filepath.Rel(htmlDir, localAbs)
						if err != nil {
							// fallback to root-relative
							rel = filepath.ToSlash(u.Path)
						}
						n.Attr[i].Val = filepath.ToSlash(rel)
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

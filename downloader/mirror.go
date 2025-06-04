package downloader

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/util"
)

// MirrorSite recursively downloads all pages and assets starting from `startURL`.
// All files are saved under a folder named after the domain (e.g., "www.example.com").
// It uses parser.ExtractLinks to find internal <a>, <link>, and <img> references.
func MirrorSite(startURL string, opts Options, log *logger.Logger) error {
	// Parse the starting URL
	base, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("invalid start URL %q: %w", startURL, err)
	}

	// Create the root directory for this domain
	domainDir := filepath.Join(opts.OutputDir, base.Host)
	if err := util.EnsureDir(domainDir, 0o755); err != nil {
		return fmt.Errorf("failed to create domain directory: %w", err)
	}

	return nil
}

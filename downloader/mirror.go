package downloader

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"

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

	// visited keeps track of which URLs we've already enqueued/downloaded
	visited := make(map[string]bool)
	// mu protects visited and queueSlice
	var mu sync.Mutex

	// queueSlice implements a simple FIFO queue
	queueSlice := []string{startURL}
	visited[startURL] = true

	for len(queueSlice) > 0 {
		// Pop front
		mu.Lock()
		currentURL := queueSlice[0]
		queueSlice = queueSlice[1:]
		mu.Unlock()

		// Download the current URL
		log.Start(currentURL, time.Now())
		resp, err := http.Get(currentURL)
		if err != nil {
			log.Error(fmt.Errorf("failed HTTP GET %s: %w", currentURL, err))
			continue
		}

	}

	return nil
}

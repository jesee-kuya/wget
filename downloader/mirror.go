package downloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/parser"
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

	visited := make(map[string]bool)
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

		urlParsed, err := url.Parse(currentURL)
		if err != nil {
			log.Error(fmt.Errorf("failed to parse URL %s: %w", currentURL, err))
			continue
		}

		if util.ShouldReject(urlParsed.Path, opts.Reject) || util.ShouldExclude(urlParsed.Path, opts.Exclude) {
			continue
		}

		// Download the current URL
		log.Start(currentURL, time.Now())
		resp, err := http.Get(currentURL)
		if err != nil {
			log.Error(fmt.Errorf("failed HTTP GET %s: %w", currentURL, err))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Error(fmt.Errorf("bad status for %s: %s", currentURL, resp.Status))
			resp.Body.Close()
			continue
		}

		// Determine where to save this file in the domain directory
		saveDir, err := util.CreateURLDirectories(currentURL, opts.OutputDir)
		if err != nil {
			log.Error(fmt.Errorf("failed to create folders for %s: %w", currentURL, err))
			resp.Body.Close()
			continue
		}

		// Derive the filename from the URL path
		filename := util.ExtractFilenameFromURL(currentURL)
		outputPath := filepath.Join(saveDir, filename)
		log.SavingTo(outputPath)

		// Read the entire response body into memory (so we can both save and parse it if HTML)
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			log.Error(fmt.Errorf("failed to read body %s: %w", currentURL, readErr))
			continue
		}

		// Write the body to a file
		fileErr := os.WriteFile(outputPath, bodyBytes, 0o644)
		if fileErr != nil {
			log.Error(fmt.Errorf("failed to write file %s: %w", outputPath, fileErr))
			continue
		}
		log.ContentInfo(int64(len(bodyBytes)))
		log.Progress(int64(len(bodyBytes)), int64(len(bodyBytes)), 0, 0)
		log.Done(time.Now(), currentURL)

		// If content-type is HTML, parse for additional links
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			// Create a reader over the saved body bytes
			reader := bytes.NewReader(bodyBytes)

			// Extract internal links
			foundLinks, parseErr := parser.ExtractLinks(base, reader)
			if parseErr != nil {
				log.Error(fmt.Errorf("failed to parse HTML %s: %w", currentURL, parseErr))
				continue
			}

			// Enqueue each new link
			for _, link := range foundLinks {
				mu.Lock()
				if !visited[link] {
					visited[link] = true
					queueSlice = append(queueSlice, link)
				}
				mu.Unlock()
			}
		}
	}

	return nil
}

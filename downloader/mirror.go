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
	base, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("invalid start URL %q: %w", startURL, err)
	}

	visited := make(map[string]bool)
	var mu sync.Mutex

	queueSlice := []string{startURL}
	visited[startURL] = true

	for len(queueSlice) > 0 {
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

		saveDir, err := util.CreateURLDirectories(currentURL, opts.OutputDir)
		if err != nil {
			log.Error(fmt.Errorf("failed to create folders for %s: %w", currentURL, err))
			resp.Body.Close()
			continue
		}

		filename := util.ExtractFilenameFromURL(currentURL)
		outputPath := filepath.Join(saveDir, filename)
		log.SavingTo(outputPath)

		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			log.Error(fmt.Errorf("failed to read body %s: %w", currentURL, readErr))
			continue
		}

		fileErr := os.WriteFile(outputPath, bodyBytes, 0o644)
		if fileErr != nil {
			log.Error(fmt.Errorf("failed to write file %s: %w", outputPath, fileErr))
			continue
		}
		log.ContentInfo(int64(len(bodyBytes)))
		log.Progress(int64(len(bodyBytes)), int64(len(bodyBytes)), 0, 0)
		log.Done(time.Now(), currentURL)

		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			reader := bytes.NewReader(bodyBytes)

			foundLinks, parseErr := parser.ExtractLinks(base, reader)
			if parseErr != nil {
				log.Error(fmt.Errorf("failed to parse HTML %s: %w", currentURL, parseErr))
				continue
			}

			for _, link := range foundLinks {
				mu.Lock()
				if !visited[link] {
					visited[link] = true
					queueSlice = append(queueSlice, link)
				}
				mu.Unlock()
			}
		}

		if opts.ConvertLink && strings.HasPrefix(contentType, "text/html") {
			rewritten, err := parser.RewriteLinks(bodyBytes, urlParsed, opts.OutputDir)
			if err != nil {
				log.Error(fmt.Errorf("rewrite links failed: %w", err))
			} else {
				bodyBytes = rewritten
			}
		}

		if err := os.WriteFile(outputPath, bodyBytes, 0o644); err != nil {
			log.Error(fmt.Errorf("write failed %s: %w", outputPath, err))
			continue
		}
	}

	return nil
}

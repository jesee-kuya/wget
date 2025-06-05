package downloader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jesee-kuya/wget/logger"
)

// DownloadInput reads URLs from opt.InputFile and downloads each via DownloadFile.
// Progress/logging is delegated to the shared logger.
func DownloadInput(opt Options, log *logger.Logger) {
	urls, err := ReadURLs(opt.InputFile)
	if err != nil {
		fmt.Fprintf(log.Output, "Error reading urls: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var completedURLs []string

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			// For each URL, we reuse DownloadFile to handle fetching, buffering, progress, etc.
			err := DownloadFile(u, opt, log)
			if err != nil {
				fmt.Fprintf(log.Output, "Error downloading %s: %v\n", u, err)
				return
			}

			// Record successful completion
			mu.Lock()
			completedURLs = append(completedURLs, u)
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	// Print summary
	fmt.Fprintf(log.Output, "Download finished: %v\n", completedURLs)
}

// Function to read URLs from the file
func ReadURLs(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", fileName, err)
	}

	// Return error if no URLs were found
	if len(urls) == 0 {
		return nil, fmt.Errorf("no valid URLs found in %s", fileName)
	}

	return urls, nil
}

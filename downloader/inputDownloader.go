package downloader

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/util"
)

func DownloadInput(opt Options, log *logger.Logger) {
	urls, err := ReadURLs(opt.InputFile)
	if err != nil {
		fmt.Fprintf(log.Output, "Error reading urls: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var completedURLs []string
	var contentSizes []int64
	var fileNames []string

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(log.Output, "Error downloading %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Fprintf(log.Output, "Error downloading %s: HTTP status %s\n", url, resp.Status)
				return
			}

			fileName := path.Base(resp.Request.URL.Path)
			if fileName == "" || fileName == "/" {
				fileName = "downloaded_file"
			}

			savePath, err := util.ProcessDirectoryPath(opt.OutputDir, true, 0o755)
			if err != nil {
				fmt.Fprintf(log.Output, "failed to process output directory: %v\n", err)
				return
			}

			fullPath := filepath.Join(savePath, fileName)
			if err := os.MkdirAll(savePath, 0o755); err != nil {
				fmt.Fprintf(log.Output, "failed to create directory %s: %v\n", savePath, err)
				return
			}

			out, err := os.Create(fullPath)
			if err != nil {
				fmt.Fprintf(log.Output, "Error creating file %s: %v\n", fileName, err)
				return
			}
			defer out.Close()

			contentLen := resp.ContentLength
			var written int64
			buf := make([]byte, 32*1024) // 32KB buffer
			start := time.Now()

			for {
				nr, er := resp.Body.Read(buf)
				if nr > 0 {
					nw, ew := out.Write(buf[:nr])
					if ew != nil {
						fmt.Fprintf(log.Output, "Error writing to %s: %v\n", fileName, ew)
						return
					}
					if nw != nr {
						fmt.Fprintf(log.Output, "Short write to %s\n", fileName)
						return
					}
					written += int64(nw)

					elapsed := time.Since(start).Seconds()
					speed := float64(written) / elapsed
					eta := time.Duration(float64(contentLen-written)/speed) * time.Second

					log.Progress(written, contentLen, speed, eta)
				}

				if er != nil {
					if er != io.EOF {
						fmt.Fprintf(log.Output, "Read error for %s: %v\n", fileName, er)
					}
					break
				}
			}

			mu.Lock()
			contentSizes = append(contentSizes, written)
			fileNames = append(fileNames, fileName)
			completedURLs = append(completedURLs, url)
			mu.Unlock()
		}(u)
	}

	wg.Wait()

	fmt.Fprintf(log.Output, "content size: %v\n", contentSizes)
	for _, name := range fileNames {
		fmt.Fprintf(log.Output, "finished %s\n", name)
	}
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

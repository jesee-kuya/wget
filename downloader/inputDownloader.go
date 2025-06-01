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

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/util"
)

func DownloadInput(opt Options, log *logger.Logger){
	
    urls, err := ReadURLs(opt.InputFile)
    if err != nil{
        fmt.Fprintf(log.Output,"Error reading urls: %v\n", err)
        return
    }

    var wg sync.WaitGroup
	var mu sync.Mutex
	var completedURLs []string
	var contentSizes []int64
	var fileNames []string

    for _, u := range urls{
        wg.Add(1)
        go func (url string){
            defer wg.Done()
    
            resp, err := http.Get(url)
            if err != nil{
                fmt.Fprintf(log.Output,"Error downloading %s: %v\n", url, err)
                return
            }
            defer resp.Body.Close()
    
            // Check if the response status is OK
            if resp.StatusCode != http.StatusOK {
                fmt.Fprintf(log.Output,"Error downloading %s: HTTP status %s\n", url, resp.Status)
                return
            }
    
    
            // Extract file name from URL
            fileName := path.Base(resp.Request.URL.Path)
            if fileName == "" || fileName == "/" {
                fileName = "downloaded_file"
            }
    

			savePath, err := util.ProcessDirectoryPath(opt.OutputDir, true, 0o755)
			if err != nil {
				fmt.Fprintf(log.Output,"failed to process output directory: %w", err)
				return
			}
    
            fullPath := filepath.Join(savePath, fileName)
    
            if err := os.MkdirAll(savePath, 0o755); err != nil {
                fmt.Fprintf(log.Output,"failed to create directory %s: %v\n", savePath, err)
                return
            }
    
    
            // Create the output file
            out, err := os.Create(fullPath)
            if err != nil {
                fmt.Fprintf(log.Output,"Error creating file %s: %v\n", fileName, err)
                return
            }
            defer out.Close()
    
            // Copy the response body to the file and track size
            size, err := io.Copy(out, resp.Body)
            if err != nil {
                fmt.Fprintf(log.Output,"Error writing to file %s: %v\n", fileName, err)
                return
            }
    
            // Safely append to completed URL, centent size and filename
            mu.Lock()
            contentSizes = append(contentSizes, size)
            fileNames = append(fileNames, fileName)
            completedURLs = append(completedURLs, url)
            mu.Unlock()
        }(u)
    }

    wg.Wait()

    // Print content size and completion message
    fmt.Fprintf(log.Output, "content size: %v\n", contentSizes)

    for _, name := range fileNames{
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

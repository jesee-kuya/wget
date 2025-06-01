package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

func downloadFile(url string, wg *sync.WaitGroup, mu *sync.Mutex, completed *[]string){
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil{
		fmt.Printf("Error downloading %s: %v\n", url, err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Error downloading %s: HTTP status %s\n", url, resp.Status)
        return
    }


	// Extract file name from URL
    fileName := path.Base(resp.Request.URL.Path)
    if fileName == "" || fileName == "/" {
        fileName = "downloaded_file"
    }

	// Create the output file
    out, err := os.Create(fileName)
    if err != nil {
        fmt.Printf("Error creating file %s: %v\n", fileName, err)
        return
    }
    defer out.Close()

	// Copy the response body to the file and track size
    size, err := io.Copy(out, resp.Body)
    if err != nil {
        fmt.Printf("Error writing to file %s: %v\n", fileName, err)
        return
    }

	// Print content size and completion message
    fmt.Printf("content size: %d\n", size)
    fmt.Printf("finished %s\n", fileName)

	// Safely append to completed URLs
    mu.Lock()
    *completed = append(*completed, url)
    mu.Unlock()
}
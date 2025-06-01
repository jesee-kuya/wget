package util

import (
	"fmt"
	"net/http"
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


}
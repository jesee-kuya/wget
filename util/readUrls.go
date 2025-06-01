package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to read URLs from the file
func readURLs(fileName string) ([]string, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
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

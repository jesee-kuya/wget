package util

import (
    "bufio"
    "os"
    "strings"
)

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

    return urls, scanner.Err()
}

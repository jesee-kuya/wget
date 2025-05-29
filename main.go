package main

import (
	"fmt"
	"os"

	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <url>")
		os.Exit(1)
	}
	url := os.Args[1]

	// Set up logger
	log := logger.NewLogger(os.Stdout)

	// Use default options (no output name or directory)
	opts := downloader.Options{
		OutputName: "",
		OutputDir:  "",
	}

	// Run the download
	err := downloader.DownloadFile(url, opts, log)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

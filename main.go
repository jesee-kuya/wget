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
		return
	}
	url := os.Args[1]

	log := logger.NewLogger(os.Stdout)

	opts := downloader.Options{
		OutputName: "",
		OutputDir:  "",
	}

	err := downloader.DownloadFile(url, opts, log)
	if err != nil {
		log.Error(err)
		return
	}
}

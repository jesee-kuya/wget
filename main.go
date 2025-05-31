package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/util"
)

func main() {
	background := flag.Bool("B", false, "Download in background and log output to wget-log")
	output := flag.String("O", "", "Specify file name")
	outputDir := flag.String("P", "", "Specify directory to save the file")
	rateLimit := flag.String("rate-limit", "", "Limit download speed (e.g., 100k, 1M)")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: wget [options] <URL>")
		return
	}

	url := args[0]

	parsedRate, err := util.ParseRateLimit(*rateLimit)
	if err != nil {
		fmt.Println("Error parsing rate limit:", err)
		return
	}

	opts := downloader.Options{
		OutputName: *output,
		OutputDir:  *outputDir,
		RateLimit:  parsedRate,
	}

	if *background {
		fmt.Println("Output will be written to \"wget-log\".")

		logFile, err := os.Create("wget-log")
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		defer logFile.Close()

		fileLogger := logger.NewLogger(logFile)

		// Perform the download using our downloader
		err = downloader.DownloadFile(url, opts, fileLogger)
		if err != nil {
			fmt.Fprintf(logFile, "Download failed: %v\n", err)
			return
		}
		return
	}

	log := logger.NewLogger(os.Stdout)

	err = downloader.DownloadFile(url, opts, log)
	if err != nil {
		log.Error(err)
		return
	}
}

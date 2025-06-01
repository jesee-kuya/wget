package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

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

		// If already running in background, skip re-exec
		if os.Getenv("WGET_BACKGROUND") != "1" {
			// Re-exec the same binary in background
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println("Error finding executable path:", err)
				return
			}

			cmd := exec.Command(execPath, os.Args[1:]...)
			cmd.Env = append(os.Environ(), "WGET_BACKGROUND=1")

			// Redirect output to log file
			logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			defer logFile.Close()

			cmd.Stdout = logFile
			cmd.Stderr = logFile
			cmd.Stdin = nil // disconnect input

			// Start the process and exit the parent
			if err := cmd.Start(); err != nil {
				fmt.Println("Failed to start background process:", err)
				return
			}
			return
		}

		// Actual background logic starts here
		logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		defer logFile.Close()

		fileLogger := logger.NewLogger(logFile)

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

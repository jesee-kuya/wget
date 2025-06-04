// cmd/wget/main.go

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
	rejectList := ""
	background := flag.Bool("B", false, "Download in background and log output to wget-log")
	output := flag.String("O", "", "Specify file name")
	inputFile := flag.String("i", "", "Input file containing URLs (one per line)")
	outputDir := flag.String("P", "", "Specify directory to save the file")
	rateLimit := flag.String("rate-limit", "", "Limit download speed (e.g., 100k, 1M)")
	mirror := flag.Bool("mirror", false, "Mirror the entire website starting from the given URL")
	rejectShort := flag.String("R", "", "Comma-separated suffixes to reject (e.g. jpg,gif)")
	reject := flag.String("reject", "", "Comma-separated suffixes to reject (e.g. jpg,gif)")

	excludeList := flag.String("X", "", "Comma-separated directories to exclude (e.g. /js,/assets)")

	flag.Parse()
	args := flag.Args()

	if *rejectShort != "" {
		rejectList = *rejectShort
	} else {
		rejectList = *reject
	}

	var urlArg string
	if *inputFile == "" {
		if len(args) == 0 {
			fmt.Println("Usage: wget [options] <URL>")
			return
		}
		urlArg = args[0]
	} else {
		urlArg = ""
	}

	parsedRate, err := util.ParseRateLimit(*rateLimit)
	if err != nil {
		fmt.Println("Error parsing rate limit:", err)
		return
	}

	opts := downloader.Options{
		OutputName:  *output,
		OutputDir:   *outputDir,
		InputFile:   *inputFile,
		RateLimit:   parsedRate,
		RunInBg:     *background,
		LogFilePath: "wget-log",
		Reject:      util.SplitAndTrim(rejectList, ","),
		Exclude:     util.SplitAndTrim(*excludeList, ","),
	}

	if *background {
		fmt.Println("Output will be written to \"wget-log\".")

		if os.Getenv("WGET_BACKGROUND") != "1" {
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println("Error finding executable path:", err)
				return
			}

			cmd := exec.Command(execPath, os.Args[1:]...)
			cmd.Env = append(os.Environ(), "WGET_BACKGROUND=1")

			logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			defer logFile.Close()

			cmd.Stdout = logFile
			cmd.Stderr = logFile
			cmd.Stdin = nil

			if err := cmd.Start(); err != nil {
				fmt.Println("Failed to start background process:", err)
				return
			}
			return
		}

		logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		defer logFile.Close()

		fileLogger := logger.NewLogger(logFile)

		if *inputFile != "" {
			downloader.DownloadInput(opts, fileLogger)
		} else if *mirror {
			err := downloader.MirrorSite(urlArg, opts, fileLogger)
			if err != nil {
				fmt.Fprintf(logFile, "Mirror failed: %v\n", err)
			}
		} else {
			err := downloader.DownloadFile(urlArg, opts, fileLogger)
			if err != nil {
				fmt.Fprintf(logFile, "Download failed: %v\n", err)
			}
		}
		return
	}

	stdoutLogger := logger.NewLogger(os.Stdout)

	if *inputFile != "" {
		downloader.DownloadInput(opts, stdoutLogger)
	} else if *mirror {
		err := downloader.MirrorSite(urlArg, opts, stdoutLogger)
		if err != nil {
			stdoutLogger.Error(err)
		}
	} else {
		err := downloader.DownloadFile(urlArg, opts, stdoutLogger)
		if err != nil {
			stdoutLogger.Error(err)
		}
	}
}

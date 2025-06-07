package worker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/logger"
)

func RunInBackground(opts downloader.Options, urlArg string) {
	fmt.Println("Output will be written to \"wget-log\".")

	if os.Getenv("WGET_BACKGROUND") != "1" {
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("Error finding executable path:", err)
			return
		}

		cmd := exec.Command(execPath, os.Args[1:]...)
		cmd.Env = append(os.Environ(), "WGET_BACKGROUND=1")

		logFile, err := os.OpenFile(opts.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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
		}
		return
	}

	logFile, err := os.OpenFile(opts.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	fileLogger := logger.NewLogger(logFile)
	Execute(opts, urlArg, *fileLogger)
}

package main

import (
	"os"
	"wget/logger"
)

func main() {
	// Create logger for stdout
	stdoutLogger := logger.NewLogger(os.Stdout)

	// Create logger for file
	file, err := os.OpenFile("wget-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		stdoutLogger.Error("Failed to open log file", err)
		return
	}
	defer file.Close()
	fileLogger := logger.NewLogger(file)

	// Example usage
	stdoutLogger.Info("Starting download")
	stdoutLogger.Progress(1024, 4096, "https://example.com/file maman")
	stdoutLogger.Error("Download failed", err)

	fileLogger.Info("Starting background download")
	fileLogger.Progress(2048, 4096, "https://example.com/file2")
}

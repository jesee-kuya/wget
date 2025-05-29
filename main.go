package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/logger"
)

func main() {
	background := flag.Bool("B", false, "Download in background and log output to wget-log")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: wget [options] <URL>")
		return
	}

	url := args[0]

	if *background {
		fmt.Println("Output will be written to \"wget-log\".")
		cmdArgs := append([]string{"wget"}, args...)
		cmd := exec.Command(os.Args[0], cmdArgs...)

		logFile, err := os.Create("wget-log")
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		defer logFile.Close()
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		err = cmd.Start()
		if err != nil {
			fmt.Println("Failed to start background download:", err)
			return
		}
		return
	}

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

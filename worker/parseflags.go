package worker

import (
	"flag"
	"fmt"
	"os"

	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/util"
)

func ParseFlags() (downloader.Options, string, bool) {
	rejectList := ""
	excludeList := ""
	background := flag.Bool("B", false, "Download in background and log output to wget-log")
	output := flag.String("O", "", "Specify file name")
	inputFile := flag.String("i", "", "Input file containing URLs (one per line)")
	outputDir := flag.String("P", "", "Specify directory to save the file")
	rateLimit := flag.String("rate-limit", "", "Limit download speed (e.g., 100k, 1M)")
	mirror := flag.Bool("mirror", false, "Mirror the entire website starting from the given URL")
	reject := flag.String("reject", "", "Comma-separated suffixes to reject (e.g. jpg,gif)")
	rejectShort := flag.String("R", "", "Comma-separated suffixes to reject (e.g. jpg,gif)")
	exclude := flag.String("exclude", "", "Comma-separated directories to exclude (e.g. /js,/assets)")
	excludeShort := flag.String("X", "", "Comma-separated directories to exclude (e.g. /js,/assets)")
	convertLinks := flag.Bool("convert-links", false, "convert the links in the downloaded files so that they can be viewed offline")

	flag.Parse()
	args := flag.Args()

	if *rejectShort != "" {
		rejectList = *rejectShort
	} else {
		rejectList = *reject
	}

	if *excludeShort != "" {
		excludeList = *excludeShort
	} else if *exclude != "" {
		excludeList = *exclude
	}

	var urlArg string
	if *inputFile == "" {
		if len(args) == 0 {
			fmt.Println("Usage: wget [options] <URL>")
			os.Exit(1)
		}
		urlArg = args[0]
	} else {
		urlArg = ""
	}

	parsedRate, err := util.ParseRateLimit(*rateLimit)
	if err != nil {
		fmt.Println("Error parsing rate limit:", err)
		os.Exit(1)
	}

	opts := downloader.Options{
		OutputName:  *output,
		OutputDir:   *outputDir,
		InputFile:   *inputFile,
		RateLimit:   parsedRate,
		RunInBg:     *background,
		LogFilePath: "wget-log",
		Reject:      util.SplitAndTrim(rejectList, ","),
		Exclude:     util.SplitAndTrim(excludeList, ","),
		ConvertLink: *convertLinks,
		Mirror:      *mirror,
	}

	return opts, urlArg, *background
}

package main

import (
	"os"

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/worker"
)

func main() {
	opts, urlArg, runBg := worker.ParseFlags()

	if runBg {
		worker.RunInBackground(opts, urlArg)
		return
	}

	logger := logger.NewLogger(os.Stdout)
	worker.Execute(opts, urlArg, *logger)
}

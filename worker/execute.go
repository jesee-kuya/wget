package worker

import (
	"github.com/jesee-kuya/wget/downloader"
	"github.com/jesee-kuya/wget/logger"
)

func Execute(opts downloader.Options, urlArg string, log logger.Logger) {
	if opts.InputFile != "" {
		downloader.DownloadInput(opts, &log)
	} else if opts.Mirror {
		err := downloader.MirrorSite(urlArg, opts, &log)
		if err != nil {
			log.Error(err)
		}
	} else {
		err := downloader.DownloadFile(urlArg, opts, &log)
		if err != nil {
			log.Error(err)
		}
	}
}

package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jesee-kuya/wget/logger"
	"github.com/jesee-kuya/wget/util"
)

// DownloadFile downloads a file from the specified URL and saves it to the output directory.
func DownloadFile(url string, opts Options, log *logger.Logger) error {
	startTime := time.Now()
	log.Start(url, startTime)

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad status from: %s, status code: %d", url, resp.StatusCode)
		log.Error(err)
		return err
	}

	log.Status(resp.StatusCode)
	log.ContentInfo(resp.ContentLength)

	// Determine output path
	filename := opts.OutputName
	if filename == "" {
		filename = util.ExtractFilenameFromURL(url)
	}
	outputPath := filepath.Join(util.FallbackDir(opts.OutputDir), filename)
	log.SavingTo(outputPath)

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer outFile.Close()

	const bufSize = 32 * 1024
	buf := make([]byte, bufSize)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var written int64
	start := time.Now()

	done := make(chan error, 1)

	go func() {
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				nw, writeErr := outFile.Write(buf[:n])
				if writeErr != nil {
					done <- writeErr
					return
				}
				if nw != n {
					done <- io.ErrShortWrite
					return
				}
				written += int64(nw)
			}

			if readErr != nil {
				if readErr == io.EOF {
					done <- nil
				} else {
					done <- readErr
				}
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			if elapsed > 0 {
				speed := float64(written) / elapsed
				eta := time.Duration(float64(resp.ContentLength-written)/speed) * time.Second
				log.Progress(written, resp.ContentLength, speed, eta)
			}
		case err := <-done:
			if err != nil {
				log.Error(err)
				return err
			}
			log.Progress(written, resp.ContentLength, float64(written)/time.Since(start).Seconds(), 0)
			log.Done(time.Now(), url)
			return nil
		}
	}
}

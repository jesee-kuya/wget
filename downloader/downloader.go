package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jesee-kuya/wget/logger"
)

func DownloadFile(url string, opts Options, log *logger.Logger) error {
	startTime := time.Now()
	log.Start(url, startTime)

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	log.Status(resp.StatusCode)
	log.ContentInfo(resp.ContentLength)

	// Build file path
	filename := opts.OutputName
	if filename == "" {
		filename = utils.ExtractFilenameFromURL(url)
	}
	outputPath := filepath.Join(utils.FallbackDir(opts.OutputDir), filename)
	log.SavingTo(outputPath)

	// Create file
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer outFile.Close()

	// Progress tracking
	const bufSize = 32 * 1024
	buf := make([]byte, bufSize)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var (
		written int64
		start   = time.Now()
	)

	// Define a multi-reader that logs progress concurrently
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			nw, writeErr := outFile.Write(buf[:n])
			if writeErr != nil {
				log.Error(writeErr)
				return writeErr
			}
			if nw != n {
				err := fmt.Errorf("short write: expected %d, wrote %d", n, nw)
				log.Error(err)
				return io.ErrShortWrite
			}
			written += int64(nw)
		}

		select {
		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			if elapsed > 0 {
				speed := float64(written) / elapsed
				eta := time.Duration(float64(resp.ContentLength-written)/speed) * time.Second
				log.Progress(written, resp.ContentLength, speed, eta)
			}
		default:
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			log.Error(readErr)
			return readErr
		}
	}

	log.Done(time.Now(), url)
	return nil
}

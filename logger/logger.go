package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jesee-kuya/wget/util"
)

type Logger struct {
	Output io.Writer
}

// NewLogger creates a new Logger instance with the specified writer
func NewLogger(output io.Writer) *Logger {
	return &Logger{Output: output}
}

// Done logs the completion of a download operation.
func (l *Logger) Done(timestamp time.Time, url string) {
	fmt.Fprintf(l.Output, "\nDownloaded [%s]\n", url)
	fmt.Fprintf(l.Output, "finished at %s\n", timestamp.Format("2006-01-02 15:04:05"))
}

// Start logs the start of a process with a timestamp.
func (l *Logger) Start(url string, timestamp time.Time) {
	fmt.Fprintf(l.Output, "start at %s\n", timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(l.Output, "sending request, awaiting response... ")
}

// SavingTo logs the path where the file is being saved.
func (l *Logger) SavingTo(path string) {
	fmt.Fprintf(l.Output, "saving file to: %s\n", path)
}

// ContentInfo logs the size of the content being downloaded.
func (l *Logger) ContentInfo(size int64) {
	readable := util.ContentSize(size)
	fmt.Fprintf(l.Output, "content size: %d [~%s]\n", size, readable)
}

// Output the status code of the process
func (l *Logger) Status(code int) {
	status := http.StatusText(code)
	if status == "" {
		status = "Unknown Status"
	}
	fmt.Fprintf(l.Output, "sending request, awaiting response... status %d %s\n", code, status)
}

// Output the progress of download
func (l *Logger) Progress(written, total int64, speed float64, eta time.Duration) {
	const barWidth = 30

	toKiB := func(b int64) float64 { return float64(b) / 1024.0 }
	writtenKiB := toKiB(written)
	speedStr := util.FormatSpeed(speed)

	if total <= 0 {
		progressIndex := int(written/10240) % barWidth
		bar := make([]rune, barWidth)
		for i := range bar {
			if i == progressIndex {
				bar[i] = '>'
			} else {
				bar[i] = ' '
			}
		}
		progressLine := fmt.Sprintf(
			"%.2f KiB / ??.?? KiB [%s]   ??%% %s ETA: ?",
			writtenKiB,
			string(bar),
			speedStr,
		)

		if l.Output == os.Stdout {
			fmt.Fprintf(l.Output, "\r%s", progressLine)
		} else {
			fmt.Fprintln(l.Output, progressLine)
		}
		return
	}

	totalKiB := toKiB(total)
	percent := float64(written) / float64(total)
	doneBars := int(percent * float64(barWidth))
	if doneBars > barWidth {
		doneBars = barWidth
	}
	if doneBars < 0 {
		doneBars = 0
	}
	remainingBars := barWidth - doneBars

	progressLine := fmt.Sprintf(
		"%.2f KiB / %.2f KiB [%s%s] %6.2f%% %s %s",
		writtenKiB,
		totalKiB,
		strings.Repeat("=", doneBars),
		strings.Repeat(" ", remainingBars),
		percent*100,
		speedStr,
		util.FormatETA(eta),
	)

	if l.Output == os.Stdout {
		fmt.Fprintf(l.Output, "\r%s", progressLine)
		if written == total {
			fmt.Fprintln(l.Output)
		}
	} else {
		fmt.Fprintln(l.Output, progressLine)
	}
}

// Error logs an error message.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.Output, "error: %v\n", err)
}

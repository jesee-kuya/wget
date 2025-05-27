package logger

import (
	"fmt"
	"io"
	"time"
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
	sizeMB := float64(size) / (1024 * 1024)
	fmt.Fprintf(l.Output, "content size: %d [~%.2fMB]\n", size, sizeMB)
}

// Error logs an error message.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.Output, "error: %v\n", err)
}

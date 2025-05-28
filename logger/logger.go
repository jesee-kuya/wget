package logger

import (
	"fmt"
	"io"
	"net/http"
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
	const (
		MB = 1024.0 * 1024.0
		GB = 1024.0 * 1024.0 * 1024.0
	)

	var readable string
	if float64(size) >= GB {
		readable = fmt.Sprintf("~%.2fGB", float64(size)/GB)
	} else {
		readable = fmt.Sprintf("~%.2fMB", float64(size)/MB)
	}

	fmt.Fprintf(l.Output, "content size: %d [%s]\n", size, readable)
}

//Output the status code of the process
func (l *Logger) Status(code int){
	status :=http.StatusText(code)
	if status == ""{
		status = "Unknown Status"
	}
	fmt.Fprintf(l.Output, "sending request, awaiting response... status %d %s\n",code , status)
}

// Error logs an error message.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.Output, "error: %v\n", err)
}

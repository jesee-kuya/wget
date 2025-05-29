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
	readable := contentSize(size)
	fmt.Fprintf(l.Output, "content size: %d [~%s]\n", size, readable)
}

//ContentSize converts size of the content into GB or MB
func contentSize(size int64) string{
	const (
		KiB = 1024.0
		MiB = 1024.0 * KiB
		GiB = 1024.0 * MiB
	)
	
	if float64(size) >= GiB {
		return fmt.Sprintf("%.2fGiB", float64(size)/GiB)
	} else if float64(size) >= MiB{
		return fmt.Sprintf("%.2fMiB", float64(size)/MiB)
	}else if float64(size) >= KiB{
		return fmt.Sprintf("%.2fKiB", float64(size)/KiB)
	}
	
	return fmt.Sprint("%.2fB", float64(size))
}

//Output the status code of the process
func (l *Logger) Status(code int){
	status :=http.StatusText(code)
	if status == ""{
		status = "Unknown Status"
	}
	fmt.Fprintf(l.Output, "sending request, awaiting response... status %d %s\n",code , status)
}

//Output the progress of download
func (l *Logger) Progress(written, total int64, speed float64, eta time.Duration){
	if written < 0{
		written = 0;
	}

	if speed < 0{
		speed = 0;
	}

	percentage := (float64(written)/float64(total))*100.0
}

func etaFormat(eta time.Duration) string{
	if eta < time.Second{
		return "0s"
	}

	seconds := int64(eta.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	return (time.Duration(seconds) * time.Second).String()
}

// Error logs an error message.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.Output, "error: %v\n", err)
}

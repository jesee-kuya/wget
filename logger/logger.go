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
	var status string
    switch code {
    case 100:
        status = "100 Continue"
    case 101:
        status = "101 Switching Protocols"
    case 102:
        status = "102 Processing"
    case 103:
        status = "103 Early Hints"
    case 200:
        status = "200 OK"
    case 201:
        status = "201 Created"
    case 202:
        status = "202 Accepted"
    case 204:
        status = "204 No Content"
    case 206:
        status = "206 Partial Content"
    case 300:
        status = "300 Multiple Choices"
    case 301:
        status = "301 Moved Permanently"
    case 302:
        status = "302 Found"
    case 303:
        status = "303 See Other"
    case 304:
        status = "304 Not Modified"
    case 307:
        status = "307 Temporary Redirect"
    case 308:
        status = "308 Permanent Redirect"
    case 400:
        status = "400 Bad Request"
    case 401:
        status = "401 Unauthorized"
    case 403:
        status = "403 Forbidden"
    case 404:
        status = "404 Not Found"
    case 405:
        status = "405 Method Not Allowed"
    case 408:
        status = "408 Request Timeout"
    case 409:
        status = "409 Conflict"
    case 410:
        status = "410 Gone"
    case 429:
        status = "429 Too Many Requests"
    case 500:
        status = "500 Internal Server Error"
    case 501:
        status = "501 Not Implemented"
    case 502:
        status = "502 Bad Gateway"
    case 503:
        status = "503 Service Unavailable"
    case 504:
        status = "504 Gateway Timeout"
    case 505:
        status = "505 HTTP Version Not Supported"
    default:
        status = fmt.Sprintf("%d Unknown Status", code)
    }
	fmt.Fprintf(l.Output, "sending request, awaiting response... status %s\n", status)
}

// Error logs an error message.
func (l *Logger) Error(err error) {
	fmt.Fprintf(l.Output, "error: %v\n", err)
}

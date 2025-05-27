package logger

import (
	"fmt"
	"time"
)

// Done logs the completion of a download operation.
// It outputs the final URL and the timestamp at which the download finished,
// formatted as "yyyy-mm-dd hh:mm:ss" for consistency with wget-style logging.
func (l *Logger) Done(timestamp time.Time, url string) {
	fmt.Fprintf(l.Output, "\nDownloaded [%s]\n", url)
	fmt.Fprintf(l.Output, "finished at %s\n", timestamp.Format("2006-01-02 15:04:05"))
}

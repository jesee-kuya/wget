package logger

import (
	"fmt"
	"time"
)

type Logger struct{}

// Done logs the completion of a download operation.
// It outputs the final URL and the timestamp at which the download finished,
// formatted as "yyyy-mm-dd hh:mm:ss" for consistency with wget-style logging.

func (l *Logger) Done(timestamp time.Time, url string) {
	fmt.Printf("Downloaded [%s]\n", url)
	fmt.Printf("finished at %s\n", timestamp.Format("2006-01-02 15:04:05"))
}

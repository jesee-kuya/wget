package logger

import (
	"fmt"
	"time"
)

// Start logs the start of a process with a timestamp.
func (l *Logger) Start(url string, timestamp time.Time) {
	fmt.Fprintf(l.Output, "start at %s\n", timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(l.Output, "sending request, awaiting response... ")
}

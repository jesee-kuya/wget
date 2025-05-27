package logger

import (
	"fmt"
	"time"
)

type Logger struct{}

// Start logs the start of a process with a timestamp.
func (l *Logger) Start(url string, timestamp time.Time) {
	fmt.Printf("start at %s %s\n", timestamp.Format("2006-01-02 15:04:05"))
	// fmt.Printf("sending request, awaiting response... %s\n", url)

}

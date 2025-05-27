package logger

import (
	"io"
	"log"
)

// Logger handles download-related logging with configurable output
type Logger struct {
	writer *log.Logger // Logger instance with io.Writer
}

// NewLogger creates a new Logger instance with the specified writer
func NewLogger(w io.Writer) *Logger {
	return &Logger{
		writer: log.New(w, "WGET: ", log.LstdFlags),
	}
}

// Info logs informational messages
func (l *Logger) Info(msg string) {
	l.writer.Printf("INFO: %s", msg)
}

// Error logs error messages
func (l *Logger) Error(msg string, err error) {
	l.writer.Printf("ERROR: %s: %v", msg, err)
}

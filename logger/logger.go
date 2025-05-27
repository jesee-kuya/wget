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

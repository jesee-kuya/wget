package logger

import (
	"io"
)

type Logger struct {
	Output io.Writer
}

// NewLogger creates a new Logger instance with the specified writer
func NewLogger(output io.Writer) *Logger {
	return &Logger{Output: output}
}

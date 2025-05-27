package logger

import (
	"io"
	"log"
)

type Logger struct {
	Output io.Writer
}

// NewLogger creates a new Logger instance with the specified writer
func NewLogger(output io.Writer) *Logger {
	return &Logger{Output: output}
}

func (l *Logger) SavingTo(path string) {
	log.Printf("saving file to: %s", path)
}

package util

import (
	"fmt"
	"time"
)

// FormatETA formats the estimated time of arrival (ETA) in a human-readable format.
func FormatETA(eta time.Duration) string {
	if eta < 0 {
		return "??s"
	}
	secs := int(eta.Seconds())
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}
	mins := secs / 60
	secs = secs % 60
	return fmt.Sprintf("%dm%ds", mins, secs)
}

// FormatSpeed formats the speed in a human-readable format.
func FormatSpeed(bytesPerSec float64) string {
	if bytesPerSec >= 1024*1024 {
		return fmt.Sprintf("%.2f MiB/s", bytesPerSec/(1024*1024))
	}
	if bytesPerSec >= 1024 {
		return fmt.Sprintf("%.2f KiB/s", bytesPerSec/1024)
	}
	return fmt.Sprintf("%.2f B/s", bytesPerSec)
}

// Repeat returns a string that repeats the first character of the input string `count` times.
func Repeat(s string, count int) string {
	return fmt.Sprintf("%s", string([]byte(s)[0])*count)
}

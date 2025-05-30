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

// ContentSize converts size of the content into GB or MB
func ContentSize(size int64) string {
	const (
		KiB = 1024.0
		MiB = 1024.0 * KiB
		GiB = 1024.0 * MiB
	)

	if float64(size) >= GiB {
		return fmt.Sprintf("%.2fGiB", float64(size)/GiB)
	} else if float64(size) >= MiB {
		return fmt.Sprintf("%.2fMiB", float64(size)/MiB)
	} else if float64(size) >= KiB {
		return fmt.Sprintf("%.2fKiB", float64(size)/KiB)
	}

	return fmt.Sprintf("%.2fB", float64(size))
}

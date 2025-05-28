package util

import "strings"

// FallbackDir returns the provided directory, or "." if it's empty or whitespace.
func FallbackDir(dir string) string {
	if strings.TrimSpace(dir) == "" {
		return "."
	}
	return dir
}

package util

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// FallbackDir returns the provided directory, or "." if it's empty or whitespace.
func FallbackDir(dir string) string {

	if strings.TrimSpace(dir) == "" {
		return "."
	}

	if strings.HasPrefix(dir, "~") {
		usr, err := user.Current()
		if err == nil {
			dir = filepath.Join(usr.HomeDir, strings.TrimPrefix(dir, "~"))
		}
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		_ = os.MkdirAll(absPath, 0755)
	}
	return absPath
}

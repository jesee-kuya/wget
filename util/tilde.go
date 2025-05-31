package util

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandTilde expands tilde (~) notation in directory paths to absolute paths.
//
// This function handles two forms of tilde expansion:
//   - "~" or "~/" expands to the current user's home directory
//   - "~username" or "~username/" expands to the specified user's home directory

func ExpandTilde(dir string) (string, error) {
	if !strings.HasPrefix(dir, "~") {
		return dir, nil
	}

	// Handle current user's home directory (~/ or just ~)
	if dir == "~" || strings.HasPrefix(dir, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get current user for tilde expansion: %w", err)
		}
		if usr.HomeDir == "" {
			return "", fmt.Errorf("current user has no home directory")
		}

		if dir == "~" {
			return usr.HomeDir, nil
		}
		return filepath.Join(usr.HomeDir, dir[2:]), nil
	}

	// Handle other user's home directory (~username or ~username/)
	var username, remainingPath string
	if slashIndex := strings.Index(dir, "/"); slashIndex != -1 {
		username = dir[1:slashIndex] // Remove ~ and get username
		remainingPath = dir[slashIndex+1:]
	} else {
		username = dir[1:] // Remove ~ and get username
	}

	if username == "" {
		return "", fmt.Errorf("empty username in tilde expansion")
	}

	usr, err := user.Lookup(username)
	if err != nil {
		return "", fmt.Errorf("failed to lookup user '%s' for tilde expansion: %w", username, err)
	}
	if usr.HomeDir == "" {
		return "", fmt.Errorf("user '%s' has no home directory", username)
	}

	if remainingPath == "" {
		return usr.HomeDir, nil
	}
	return filepath.Join(usr.HomeDir, remainingPath), nil
}

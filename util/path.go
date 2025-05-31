package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir checks if a directory exists, and creates it with the specified permissions if it does not.
func EnsureDir(dir string, perm os.FileMode) error {
	if dir == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path '%s' exists but is not a directory", dir)
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory '%s': %w", dir, err)
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	return nil
}

func ResolveToAbsolute(dir string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("directory path cannot be empty")
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for '%s': %w", dir, err)
	}

	return absPath, nil
}

// ProcessDirectoryPath processes a directory path by applying a fallback for empty input,
func ProcessDirectoryPath(dir string, createIfMissing bool, perm os.FileMode) (string, error) {
	processedDir := FallbackDir(dir)

	expandedDir, err := ExpandTilde(processedDir)
	if err != nil {
		return "", fmt.Errorf("tilde expansion failed: %w", err)
	}

	absoluteDir, err := ResolveToAbsolute(expandedDir)
	if err != nil {
		return "", fmt.Errorf("path resolution failed: %w", err)
	}

	if createIfMissing {
		if err := EnsureDir(absoluteDir, perm); err != nil {
			return "", fmt.Errorf("directory creation failed: %w", err)
		}
	}

	return absoluteDir, nil
}

// DirectoryExists checks whether the specified path exists and is a directory.
func DirectoryExists(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

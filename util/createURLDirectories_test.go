package util

import (
	"path/filepath"
	"testing"
)

var Testcases1 = []struct {
	name string
	url  string
}{
	{"Testcase2", "https://example.com/another/path/to/resource.html"},
}

func TestCreateURLDirectories(t *testing.T) {
	for _, tc := range Testcases1 {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := CreateURLDirectories(tc.url, "")
			if err != nil {
				t.Errorf("Failed to create directories for %s: %v", tc.url, err)
			} else {
				expectedDir := filepath.Join("", "example.com", "another", "path", "to")
				if dir != expectedDir {
					t.Errorf("Expected directory %s, got %s", expectedDir, dir)
				}
			}
		})
	}
}

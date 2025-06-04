package util

import "strings"

// SplitAndTrim splits s by sep and trims whitespace from each element.
// Empty elements are skipped.
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var out []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

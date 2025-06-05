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

// ShouldReject checks if the given path ends with any of the specified suffixes.
func ShouldReject(path string, rejectList []string) bool {
	if len(rejectList) == 0 {
		return false
	}
	lowerPath := strings.ToLower(path)
	for _, suffix := range rejectList {
		suff := strings.ToLower(strings.TrimSpace(suffix))
		if suff == "" {
			continue
		}
		// Ensure suffix begins with a dot
		if !strings.HasPrefix(suff, ".") {
			suff = "." + suff
		}
		if strings.HasSuffix(lowerPath, suff) {
			return true
		}
	}
	return false
}

// ShouldExclude checks if the given path starts with any of the specified prefixes.
func ShouldExclude(path string, excludeList []string) bool {
	if len(excludeList) == 0 {
		return false
	}
	for _, prefix := range excludeList {
		pref := strings.TrimSpace(prefix)
		if pref == "" {
			continue
		}
		// Compare raw path segments (case‐sensitive)
		if strings.HasPrefix(path, pref) {
			return true
		}
	}
	return false
}

package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseRateLimit(rate string) (float64, error) {
	if rate == "" {
		return 0, nil
	}

	multiplier := 1.0
	unit := strings.ToLower(rate[len(rate)-1:])

	switch unit {
	case "k":
		multiplier = 1024.0
		rate = rate[:len(rate)-1]
	case "m":
		multiplier = 1024.0 * 1024.0
		rate = rate[:len(rate)-1]
	case "g":
		multiplier = 1024.0 * 1024.0 * 1024.0
		rate = rate[:len(rate)-1]
	default:
		return 0, fmt.Errorf("invalid rate limit unit: %s", unit)

	}

	value, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse rate limit: %w", err)
	}

	return value * multiplier, nil
}

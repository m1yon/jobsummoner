package linkedincrawler

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func parseRelativeTime(now time.Time, input string) (time.Time, error) {
	// Regex to find number and unit
	re := regexp.MustCompile(`(\d+)\s*(second|minute|hour|day)s?\b\s*ago`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		return time.Time{}, fmt.Errorf("invalid format")
	}

	// Convert number to integer
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, err
	}

	// Determine the unit of time
	var duration time.Duration
	switch matches[2] {
	case "second":
		duration = time.Second * time.Duration(value)
	case "minute":
		duration = time.Minute * time.Duration(value)
	case "hour":
		duration = time.Hour * time.Duration(value)
	case "day":
		duration = time.Hour * 24 * time.Duration(value)
	default:
		return time.Time{}, fmt.Errorf("unknown time unit")
	}

	// Calculate the time
	actualTime := now.Add(-duration)
	return actualTime, nil
}

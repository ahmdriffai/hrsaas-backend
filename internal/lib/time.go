package lib

import (
	"errors"
	"fmt"
	"time"
)

var ErrInvalidTimeFormat = errors.New("invalid time format")

// ParseTimeHHMMOrHHMMSS parses "15:04" and "15:04:05" into time.Time (UTC).
func ParseTimeHHMMOrHHMMSS(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	switch len(value) {
	case 5:
		hour, minute, ok := parseHHMM(value)
		if !ok {
			return time.Time{}, ErrInvalidTimeFormat
		}
		return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC), nil
	case 8:
		hour, minute, second, ok := parseHHMMSS(value)
		if !ok {
			return time.Time{}, ErrInvalidTimeFormat
		}
		return time.Date(0, 1, 1, hour, minute, second, 0, time.UTC), nil
	default:
		return time.Time{}, ErrInvalidTimeFormat
	}
}

func parseHHMM(value string) (hour int, minute int, ok bool) {
	if value[2] != ':' {
		return 0, 0, false
	}

	hour, ok = parseTwoDigits(value[0], value[1])
	if !ok || hour > 23 {
		return 0, 0, false
	}

	minute, ok = parseTwoDigits(value[3], value[4])
	if !ok || minute > 59 {
		return 0, 0, false
	}

	return hour, minute, true
}

func parseHHMMSS(value string) (hour int, minute int, second int, ok bool) {
	if value[2] != ':' || value[5] != ':' {
		return 0, 0, 0, false
	}

	hour, ok = parseTwoDigits(value[0], value[1])
	if !ok || hour > 23 {
		return 0, 0, 0, false
	}

	minute, ok = parseTwoDigits(value[3], value[4])
	if !ok || minute > 59 {
		return 0, 0, 0, false
	}

	second, ok = parseTwoDigits(value[6], value[7])
	if !ok || second > 59 {
		return 0, 0, 0, false
	}

	return hour, minute, second, true
}

func parseTwoDigits(a, b byte) (int, bool) {
	if a < '0' || a > '9' || b < '0' || b > '9' {
		return 0, false
	}
	return int(a-'0')*10 + int(b-'0'), true
}

func ParseDateToUnixMilli(dateStr string) (int64, error) {

	// Parse the string into a time.Time object
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 0, err
	}

	// Convert to Unix epoch
	return t.UnixMilli(), nil
}

func ParseTimeToUnixMilli(timeStr string) (int64, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, err
	}

	return t.UnixMilli(), nil
}

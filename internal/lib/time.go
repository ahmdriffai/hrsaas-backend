package lib

import (
	"errors"
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

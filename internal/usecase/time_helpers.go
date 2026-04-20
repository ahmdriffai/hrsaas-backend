package usecase

import "time"

// Centralized time source for usecases.
func nowEpoch() int64 {
	return time.Now().UnixMilli()
}

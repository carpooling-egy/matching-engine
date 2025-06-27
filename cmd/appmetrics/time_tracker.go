package appmetrics

import (
	"matching-engine/internal/collections"
	"time"
)

var Timings = collections.NewSyncMap[string, time.Duration]()

// TrackTime tracks the time taken for a specific operation identified by the key.
func TrackTime(key string, duration time.Duration) {
	if existingDuration, exists := Timings.Get(key); exists {
		Timings.Set(key, existingDuration+duration)
	} else {
		Timings.Set(key, duration)
	}
}

// ResetTimings clears all tracked timings.
func ResetTimings() {
	Timings.Clear()
}

// GetAllTimings retrieves all tracked timings.
func GetAllTimings() map[string]time.Duration {
	allTimings := make(map[string]time.Duration)
	Timings.ForEach(func(key string, duration time.Duration) error {
		allTimings[key] = duration
		return nil
	})
	return allTimings
}

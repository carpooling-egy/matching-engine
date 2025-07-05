package appmetrics

import "matching-engine/internal/collections"

var Counts = collections.NewSyncMap[string, float64]()

// IncrementCount increments the count for a specific key by 1.
func IncrementCount(key string, increment float64) {
    if existingCount, exists := Counts.Get(key); exists {
        Counts.Set(key, existingCount+increment)
    } else {
        Counts.Set(key, increment)
    }
}

func GetCount(key string) float64 {
    if count, exists := Counts.Get(key); exists {
        return count
    }
    return 0.0
}

// ResetCounts clears all tracked counts.
func ResetCounts() {
    Counts.Clear()
}

// GetAllCounts retrieves all tracked counts.
func GetAllCounts() map[string]float64 {
    allCounts := make(map[string]float64)
    Counts.ForEach(func(key string, count float64) error {
        allCounts[key] = count
        return nil
    })
    return allCounts
}
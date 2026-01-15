package filter

import "time"

// FilterResult represents the result of filtering check
type FilterResult int

const (
	Include FilterResult = iota // TrxDate is within range
	Skip                        // TrxDate is before start
	Stop                        // TrxDate is >= end, stop processing
)

// CheckRange checks if a timestamp is within the range [start, end)
// Returns Include if start <= t < end
// Returns Skip if t < start
// Returns Stop if t >= end
func CheckRange(t, start, end time.Time) FilterResult {
	if t.Before(start) {
		return Skip
	}
	if !t.Before(end) { // t >= end
		return Stop
	}
	return Include
}

// InRange returns true if start <= t < end
func InRange(t, start, end time.Time) bool {
	return !t.Before(start) && t.Before(end)
}

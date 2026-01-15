package filter

import (
	"testing"
	"time"
)

func TestCheckRange(t *testing.T) {
	start := time.Date(2025, 6, 28, 0, 0, 0, 0, time.FixedZone("WIB", 7*3600))
	end := time.Date(2025, 7, 3, 0, 0, 0, 0, time.FixedZone("WIB", 7*3600))

	tests := []struct {
		name     string
		t        time.Time
		expected FilterResult
	}{
		{
			name:     "before start",
			t:        time.Date(2025, 6, 27, 23, 59, 59, 0, time.FixedZone("WIB", 7*3600)),
			expected: Skip,
		},
		{
			name:     "exactly at start (inclusive)",
			t:        start,
			expected: Include,
		},
		{
			name:     "within range",
			t:        time.Date(2025, 6, 30, 12, 0, 0, 0, time.FixedZone("WIB", 7*3600)),
			expected: Include,
		},
		{
			name:     "one second before end",
			t:        time.Date(2025, 7, 2, 23, 59, 59, 0, time.FixedZone("WIB", 7*3600)),
			expected: Include,
		},
		{
			name:     "exactly at end (exclusive)",
			t:        end,
			expected: Stop,
		},
		{
			name:     "after end",
			t:        time.Date(2025, 7, 3, 0, 0, 1, 0, time.FixedZone("WIB", 7*3600)),
			expected: Stop,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckRange(tt.t, start, end)
			if result != tt.expected {
				t.Errorf("CheckRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	start := time.Date(2025, 6, 28, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 7, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		t        time.Time
		expected bool
	}{
		{"before start", time.Date(2025, 6, 27, 0, 0, 0, 0, time.UTC), false},
		{"at start", start, true},
		{"in range", time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC), true},
		{"at end", end, false},
		{"after end", time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InRange(tt.t, start, end)
			if result != tt.expected {
				t.Errorf("InRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

package db

import "time"

// Weight represents a weight value with a timestamp.
type Weight struct {
	Value     float64
	UpdatedAt time.Time
}

// NewWeight creates a new weight value with the current timestamp.
func NewWeight(val float64) Weight {
	return Weight{
		Value:     val,
		UpdatedAt: time.Now().UTC(),
	}
}

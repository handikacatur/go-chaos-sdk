package core

import (
	"context"
	"math/rand"
	"time"
)

// InjectLatency blocks until the duration elapses or the context is canceled.
// It returns ctx.Err() if the context is closed before the duration.
func InjectLatency(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}

	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ShouldFail returns true if a random float is less than the given rate.
// Rate must be between 0.0 and 1.0.
func ShouldFail(rate float64) bool {
	if rate <= 0 {
		return false
	}
	if rate >= 1.0 {
		return true
	}
	return rand.Float64() < rate
}

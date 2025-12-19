package core

import (
	"time"
)

// Config holds the configuration for chaos injection rules.
type Config struct {
	// Enabled acts as a master switch. If false, all chaos logic is skipped.
	Enabled bool

	// HeaderTrigger is the request header key used to activate chaos.
	// If empty, the middleware may apply chaos to all requests depending on implementation.
	HeaderTrigger string

	// Latency is the duration to block the request.
	Latency time.Duration

	// FailureRate is the probability (0.0 to 1.0) that a request fails.
	FailureRate float64
}

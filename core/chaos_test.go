package core

import (
	"context"
	"testing"
	"time"
)

func TestInjectLatency(t *testing.T) {
	// Case 1: Normal Latency (Should wait)
	t.Run("waits_for_duration", func(t *testing.T) {
		start := time.Now()
		// We use a small duration to keep tests fast
		err := InjectLatency(context.Background(), 50*time.Millisecond)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if elapsed < 50*time.Millisecond {
			t.Error("function returned too fast, latency not applied")
		}
	})

	// Case 2: Context Cancellation (Should exit immediately)
	t.Run("respects_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		start := time.Now()
		// Try to sleep for 2 seconds (much longer than timeout)
		err := InjectLatency(ctx, 2*time.Second)
		elapsed := time.Since(start)

		if err == nil {
			t.Error("expected context error, got nil")
		}
		// It should have exited around 10ms, definitely not 2s
		if elapsed > 100*time.Millisecond {
			t.Errorf("function hung for %v, did not respect context cancel", elapsed)
		}
	})
}

func TestShouldFail(t *testing.T) {
	t.Run("always_fails_at_1.0", func(t *testing.T) {
		if !ShouldFail(1.0) {
			t.Error("rate 1.0 should always return true")
		}
	})

	t.Run("never_fails_at_0.0", func(t *testing.T) {
		if ShouldFail(0.0) {
			t.Error("rate 0.0 should never return true")
		}
	})
}

package httpchaos

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/handikacatur/go-chaos-sdk/core"
)

func TestMiddleware(t *testing.T) {
	// Base mock handler that always returns 200 OK
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		cfg            core.Config
		headerKey      string
		headerVal      string
		expectedStatus int
		minDuration    time.Duration
	}{
		{
			name:           "disabled_globally",
			cfg:            core.Config{Enabled: false, FailureRate: 1.0},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no_header_trigger",
			cfg:            core.Config{Enabled: true, HeaderTrigger: "x-chaos", FailureRate: 1.0},
			headerKey:      "x-other",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "trigger_failure",
			cfg:            core.Config{Enabled: true, HeaderTrigger: "x-chaos", FailureRate: 1.0},
			headerKey:      "x-chaos",
			headerVal:      "true",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "trigger_latency",
			cfg:            core.Config{Enabled: true, HeaderTrigger: "x-chaos", Latency: 50 * time.Millisecond},
			headerKey:      "x-chaos",
			headerVal:      "true",
			expectedStatus: http.StatusOK, // It eventually succeeds
			minDuration:    50 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Create the middleware
			handler := Middleware(tt.cfg)(baseHandler)

			// 2. Create a request
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			if tt.headerKey != "" {
				req.Header.Set(tt.headerKey, tt.headerVal)
			}

			// 3. Record the response
			rec := httptest.NewRecorder()
			start := time.Now()
			handler.ServeHTTP(rec, req)
			duration := time.Since(start)

			// 4. Assertions
			if rec.Code != tt.expectedStatus {
				t.Errorf("status code: expected %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.minDuration > 0 && duration < tt.minDuration {
				t.Errorf("latency: expected > %v, got %v", tt.minDuration, duration)
			}
		})
	}
}

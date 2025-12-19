package httpchaos

import (
	"net/http"

	"github.com/handikacatur/go-chaos-sdk/core"
)

// Middleware returns a standard http.Handler wrapper that injects chaos.
// It is compatible with net/http, Chi, Gorilla Mux, and Gin (via adapter).
func Middleware(cfg core.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Fast path: if chaos is globally disabled, skip all logic.
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// If a trigger header is defined, ensure it exists.
			if cfg.HeaderTrigger != "" {
				if val := r.Header.Get(cfg.HeaderTrigger); val == "" {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Inject latency using the request context.
			// If the client disconnects (context canceled), we stop immediately.
			if cfg.Latency > 0 {
				if err := core.InjectLatency(r.Context(), cfg.Latency); err != nil {
					return
				}
			}

			// Inject a transient failure (503 Service Unavailable).
			// This mimics an overloaded server or upstream outage.
			if core.ShouldFail(cfg.FailureRate) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Chaos Injected: Service Unavailable"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

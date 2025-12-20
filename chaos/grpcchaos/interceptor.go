package grpcchaos

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/handikacatur/go-chaos-sdk/chaos"
	"github.com/handikacatur/go-chaos-sdk/core"
)

// UnaryServerInterceptor returns a unary server interceptor that conditionally injects
// latency or failures into gRPC requests based on the provided configuration.
//
// It checks for a specific trigger header (if configured) before applying chaos.
func UnaryServerInterceptor(cfg chaos.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Fast path: if chaos is globally disabled, skip all logic.
		if !cfg.Enabled {
			return handler(ctx, req)
		}

		// If a trigger header is defined, ensure it exists in the incoming metadata.
		if cfg.HeaderTrigger != "" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			if values := md.Get(cfg.HeaderTrigger); len(values) == 0 {
				return handler(ctx, req)
			}
		}

		if cfg.Latency > 0 {
			if err := core.InjectLatency(ctx, cfg.Latency); err != nil {
				return nil, err
			}
		}

		// Inject a transient failure (Unavailable) to verify client retry policies.
		if core.ShouldFail(cfg.FailureRate) {
			return nil, status.Error(codes.Unavailable, "Chaos Injected: Service Unavailable")
		}

		return handler(ctx, req)
	}
}

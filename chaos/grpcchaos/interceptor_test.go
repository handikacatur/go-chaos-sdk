package grpcchaos

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/handikacatur/go-chaos-sdk/core"
)

func TestUnaryServerInterceptor(t *testing.T) {
	// Mock handler that returns success
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: "TestService/TestMethod"}

	tests := []struct {
		name        string
		cfg         core.Config
		ctx         context.Context
		expectError bool
		expectCode  codes.Code
	}{
		{
			name: "disabled_skips_chaos",
			cfg:  core.Config{Enabled: false, FailureRate: 1.0},
			ctx:  context.Background(),
		},
		{
			name: "missing_header_skips_chaos",
			cfg:  core.Config{Enabled: true, HeaderTrigger: "x-chaos", FailureRate: 1.0},
			ctx:  context.Background(), // No metadata
		},
		{
			name: "injects_failure",
			cfg:  core.Config{Enabled: true, HeaderTrigger: "x-chaos", FailureRate: 1.0},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-chaos", "true")),
			expectError: true,
			expectCode:  codes.Unavailable,
		},
		{
			name: "injects_latency",
			cfg:  core.Config{Enabled: true, HeaderTrigger: "x-chaos", Latency: 50 * time.Millisecond},
			ctx:  metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-chaos", "true")),
			// Should succeed after waiting
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := UnaryServerInterceptor(tt.cfg)
			
			start := time.Now()
			_, err := interceptor(tt.ctx, nil, info, mockHandler)
			duration := time.Since(start)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				st, _ := status.FromError(err)
				if st.Code() != tt.expectCode {
					t.Errorf("expected code %v, got %v", tt.expectCode, st.Code())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if tt.cfg.Latency > 0 && tt.expectError == false {
				// If we expected latency injection
				if duration < tt.cfg.Latency {
					t.Errorf("expected latency > %v, got %v", tt.cfg.Latency, duration)
				}
			}
		})
	}
}

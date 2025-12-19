package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/handikacatur/go-chaos-sdk/core"
	"github.com/handikacatur/go-chaos-sdk/chaos/grpcchaos"
	pb "github.com/handikacatur/go-chaos-sdk/examples/grpc/proto"
)

// server is used to implement demo.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements demo.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 1. Define Chaos Policy
	cfg := core.Config{
		Enabled:       true,
		HeaderTrigger: "x-chaos-test",
		Latency:       3 * time.Second, // 3s delay for gRPC
		FailureRate:   0.0,
	}

	// 2. Start Listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 3. Register Interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcchaos.UnaryServerInterceptor(cfg)),
	)
	
	pb.RegisterGreeterServer(s, &server{})
	
	// Enable reflection so we can use 'grpcurl' without needing .proto files
	reflection.Register(s)

	log.Printf("🚀 gRPC Server listening at %v", lis.Addr())
	log.Println("👉 Test Chaos: grpcurl -plaintext -H 'x-chaos-test: true' localhost:50051 demo.Greeter/SayHello")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

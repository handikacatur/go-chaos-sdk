package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/handikacatur/go-chaos-sdk/chaos/httpchaos"
	"github.com/handikacatur/go-chaos-sdk/chaos"
)

func main() {
	cfg := chaos.Config{
		Enabled:       true,
		HeaderTrigger: "x-chaos-test",
		Latency:       4 * time.Second,
		FailureRate:   0.0,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "pong", "status": "ok"}`))
	})

	handler := httpchaos.Middleware(cfg)(mux)

	port := "3000"
	fmt.Printf("🚀 HTTP Server running on :%s\n", port)
	fmt.Printf("👉 Test Normal: curl localhost:%s/ping\n", port)
	fmt.Printf("👉 Test Chaos:  curl -H 'x-chaos-test: true' localhost:%s/ping\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

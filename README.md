# 💥 Go Chaos SDK

![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat&logo=go)
[![Go Report Card](https://goreportcard.com/badge/github.com/handikacatur/go-chaos-sdk)](https://goreportcard.com/report/github.com/handikacatur/go-chaos-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> **Stop shipping fragile microservices.** Test network timeouts and failures locally—without Docker, Sidecars, or Root access.

**Go Chaos SDK** is a lightweight, dependency-free library for injecting controlled failures into your Go applications. It provides a unified way to test resilience across both **gRPC** and **HTTP/REST** services.

---

## 🚀 Why use this?

Testing timeouts and error handling usually requires heavy tools like
*Chaos Mesh* or *Toxiproxy*. This is overkill for unit tests or local development.

**Go Chaos SDK** solves this by living inside your middleware.

* ✅ **Zero Infrastructure:** Just a Go import. No sidecars or agents.
* ✅ **Unified Config:** Use one config for both gRPC and HTTP.
* ✅ **Header-Driven:** Trigger failures instantly via HTTP headers (great for CI/CD).
* ✅ **Safety First:** Includes a "Kill Switch" for production safety.

---

## 📦 Installation

```bash
go get https://github.com/handikacatur/go-chaos-sdk


```

---

## 🛠️ Usage

1. Create chaos configuration.

```Go
var chaosConfig = chaos.Config{
    // The "Safety Lock". Set this to false in Production!
    Enabled:       os.Getenv("APP_ENV") != "production",

    // The secret header to trigger chaos
    HeaderTrigger: "x-chaos-test",

    // Simulate a 2-second network lag
    Latency:       2 * time.Second,

    // 10% chance of returning a 503/Unavailable error
    FailureRate:   0.1, 
}
```

2. Add to your gRPC

Works as standard Unary interceptor.

```go
import (
    "google.golang.org/grpc"
    "https://github.com/handikacatur/go-chaos-sdk/grpcchaos"
)

// Add to Server Options
s := grpc.NewServer(
    grpc.UnaryInterceptor(grpcchaos.UnaryServerInterceptor(chaosConfig)),
)
```

3. Add to your HTTP/REST API

Works with net/http, Chi, and Mux.

```go
import (
    "net/http"
    "https://github.com/handikacatur/go-chaos-sdk/httpchaos"
)

// Wrap your router
mux := http.NewServeMux()
// ... register your routes ...

// Apply Middleware
handler := httpchaos.Middleware(chaosConfig)(mux)

http.ListenAndServe(":8080", handler)

```

If your framework use `fasthttp`, you need to wrap it using wrapper.

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/adaptor/v2"
    "github.com/handikacatur/go-chaos-sdk/chaos"
    "github.com/handikacatur/go-chaos-sdk/httpchaos"
)

func main() {
    app := fiber.New()

    chaosConfig := chaos.Config{
        Enabled:       true,
        HeaderTrigger: "x-chaos-test",
        Latency:       2 * time.Second,
    }

    // Fiber cannot run standard middleware directly.
    app.Use(adaptor.HTTPMiddleware(httpchaos.Middleware(chaosConfig)))

    app.Get("/ping", func(c *fiber.Ctx) error {
        return c.SendString("pong")
    })

    app.Listen(":3000")
}
```

---

## 🥽 How To Test

Once running, your service behaves normally until you send the magic header.

### Test Latency (curl)

```bash
# Normal request (Fast)
curl localhost:8080/ping

# 💥 Chaos request (Hangs for 2s)
curl -H "x-chaos-test: true" localhost:8080/ping
```

### Test gRPC (grpcurl)

```bash
# 💥 Chaos request
grpcurl -plaintext -H "x-chaos-test: true" localhost:50051 my.Service/Method
```

---

## 🛡️ Security Best Practice

⚠️ Warning: This tool is powerful. Do not deploy it to Production without safeguards.

1. Environment Gating: Always set `Enabled: false` by default.
Only enable it if `APP_ENV=staging` or `APP_ENV=dev`.

2. Strip Headers: Configure your Nginx/Load Balancer to strip the 
`x-chaos-test header` from incoming external traffic to prevent public abuse.

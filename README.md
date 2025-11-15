# uuidify-go — Fast UUID & ULID generation for Go apps
> Official Go SDK for the [UUIDify API](https://github.com/ilkereroglu/uuidify).

[![Go Reference](https://pkg.go.dev/badge/github.com/ilkereroglu/uuidify-go.svg)](https://pkg.go.dev/github.com/ilkereroglu/uuidify-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ilkereroglu/uuidify-go)](https://goreportcard.com/report/github.com/ilkereroglu/uuidify-go)
![Release](https://img.shields.io/github/v/release/ilkereroglu/uuidify-go?color=blue)
![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
[![Go SDK Tests](https://github.com/ilkereroglu/uuidify-go/actions/workflows/test.yaml/badge.svg)](https://github.com/ilkereroglu/uuidify-go/actions/workflows/test.yaml)

Minimal, idiomatic Go client for generating UUIDv1/v4/v7 and ULID identifiers through UUIDify’s globally distributed API.

---

## Install
```bash
go get github.com/ilkereroglu/uuidify-go
```

## Usage
```go
package main

import (
    "context"
    "fmt"
    "log"

    uuidify "github.com/ilkereroglu/uuidify-go"
)

func main() {
    client, err := uuidify.NewDefaultClient()
    if err != nil {
        log.Fatal(err)
    }

    id, err := client.UUIDv4(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("UUID v4:", id)
}
```

## Examples
Concrete demos live under [`examples/`](examples):

- `go run ./examples/uuid_v4` – fetch a single UUIDv4.
- `go run ./examples/uuid_v7_batch` – batch-generate v7 identifiers.
- `go run ./examples/ulid` – obtain a ULID.
Each example uses `context.Context`, the default client, and the same error handling patterns you can reuse in your services.

## Features
- ✅ Drop-in `NewDefaultClient()` with overridable base URL, HTTP client, and User-Agent.
- ⚡️ Fetch UUIDv1/v4/v7, ULID, or batch payloads with one call.
- 🧵 Context-aware HTTP requests, perfect for microservices, CLIs, and serverless workloads.
- 🎯 Typed error system (`RequestError`, `APIError`, `DecodeError`) for clean retries and observability.
- 🧩 Generated directly from UUIDify’s OpenAPI spec, ensuring long-term compatibility.
- 🧪 Backed by Go tooling (`go test`, `go vet`, CI) and production-friendly release workflow.

## Why UUIDify
UUIDify is a latency-optimized unique identifier service built for modern Go developers. With this SDK you get:

- A highly available UUID/ULID generator distributed across regions.
- Predictable performance without maintaining your own randomness infrastructure.
- Consistent REST semantics you can test locally and promote to production seamlessly.
- OpenAPI-driven definitions to keep typed clients in sync across releases.

## Service & Documentation
- Main service repository: [github.com/ilkereroglu/uuidify](https://github.com/ilkereroglu/uuidify)
- API documentation & schema: [`api/openapi.yaml`](https://github.com/ilkereroglu/uuidify/blob/main/api/openapi.yaml)

## API Reference
- Go package reference on pkg.go.dev: [pkg.go.dev/github.com/ilkereroglu/uuidify-go](https://pkg.go.dev/github.com/ilkereroglu/uuidify-go)

## License
MIT License © [ilkereroglu](https://github.com/ilkereroglu)

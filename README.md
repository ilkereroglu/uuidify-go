# uuidify-go

![Release](https://img.shields.io/github/v/release/ilkereroglu/uuidify-go?style=for-the-badge)
![Go SDK Tests](https://github.com/ilkereroglu/uuidify-go/actions/workflows/test.yaml/badge.svg)
![OpenAPI Sync](https://github.com/ilkereroglu/uuidify-go/actions/workflows/sync-openapi.yaml/badge.svg)

Official Go SDK for [uuidify.io](https://uuidify.io) offering fast UUID/ULID generation with a minimal, idiomatic API. The package ships as a single module (`github.com/ilkereroglu/uuidify-go`) and can be embedded in CLIs, services, and serverless functions without extra dependencies.

## Features
- Simple `uuidify.NewClient()` constructor with overridable base URL, HTTP client, and User-Agent
- Support for UUID v1, v4, and v7 plus ULID generation (single or batched)
- Context-aware requests with consistent error types for transport, API status, and JSON decoding issues
- Example applications and CI pipelines to keep the SDK production-ready

## Installation
```bash
go get github.com/ilkereroglu/uuidify-go
```

## Quick Start
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ilkereroglu/uuidify-go"
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

## Usage
### Customizing the Client
```go
client, err := uuidify.NewClient(
    uuidify.DefaultBaseURL,
    uuidify.WithUserAgent("my-app/1.0"),
)
if err != nil {
    log.Fatal(err)
}
```
The client exposes:
- `UUIDv1(ctx)` / `UUIDv4(ctx)` / `UUIDv7(ctx)`
- `ULID(ctx)`
- `UUIDBatch(ctx, version, count)`
- `ULIDBatch(ctx, count)`
Each method respects the provided `context.Context` for cancellation and deadlines.

### Error Handling
Three dedicated error types make it easy to react appropriately:
- `RequestError` – network failures or invalid requests
- `APIError` – non-2xx responses, including HTTP status code and message
- `DecodeError` – JSON decoding failures

Use `errors.As` to branch on these errors:
```go
var apiErr *uuidify.APIError
if errors.As(err, &apiErr) {
    fmt.Println("UUIDify returned status", apiErr.StatusCode)
}
```

## Examples
Run the samples in the [`examples/`](examples) directory:
- [`examples/uuid_v4`](examples/uuid_v4)
- [`examples/uuid_v7_batch`](examples/uuid_v7_batch)
- [`examples/ulid`](examples/ulid)
```bash
go run ./examples/uuid_v7_batch
```

## OpenAPI & Code Generation
The `openapi/` directory is auto-synced from the upstream uuidify repository via GitHub Actions. Custom templates under `templates/` are used with [`oapi-codegen`](https://github.com/deepmap/oapi-codegen) (see `openapi/config.yaml`) to regenerate the SDK (`client.go`, `models.go`, `request.go`).

## Development
```bash
go vet ./...
go test ./...
go build ./...
```
A CI workflow (`.github/workflows/test.yaml`) runs lint, build, and tests on every PR and push. Use `go test -tags live ./...` to hit the production API when needed.

## License
This project is licensed under the [MIT License](LICENSE).

# uuidify-go — Official Go SDK for uuidify.io

## Installation
```
go get github.com/ilkereroglu/uuidify-go
```

## Import
```go
import "github.com/ilkereroglu/uuidify-go"
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
    client := uuidify.NewClient()

    id, err := client.UUIDv4(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("UUID v4:", id)
}
```

## Examples
Explore runnable samples in the `examples/` directory:
- [examples/uuid_v4](examples/uuid_v4)
- [examples/uuid_v7_batch](examples/uuid_v7_batch)
- [examples/ulid](examples/ulid)

## Supported Methods
- `UUIDv1()`
- `UUIDv4()`
- `UUIDv7()`
- `ULID()`
- `UUIDBatch(version, count)`
- `ULIDBatch(count)`

Each method accepts a `context.Context` so calls can be cancelled or timed out by the caller.

## Errors
The SDK surfaces descriptive error types to help callers react appropriately:
- `RequestError`: wraps request construction or transport failures (e.g., network issues).
- `APIError`: returned when the UUIDify service responds with a non-2xx status; includes the HTTP status code and response message.
- `DecodeError`: wraps failures that occur while decoding JSON responses.

Use `errors.As` / `errors.Is` to detect these error types and implement custom retry or fallback logic.

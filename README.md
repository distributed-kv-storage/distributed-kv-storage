# Distributed KV Storage

Distributed key-value storage written in Go. The project consists of separate
`kv-engine`, `repl`, `rpc-proxy`, `kvctl`, and `clusterstat` binaries.

## API contracts

- `api/proto` contains source contracts
- `api/gen/go` contains generated messages and gRPC stubs
- internal Raft RPCs belong to `repl`

## KV API v1

[`kv.proto`](api/proto/kv/v1/kv.proto) defines `Set(key, value, ttl)` and
`Get(key)`. Keys and values are byte sequences; keys must not be empty.

TTL is expressed in seconds:

- `0` means no expiration
- values up to 30 days are relative
- larger values are Unix timestamps

Expired keys are returned as `NotFound`. Other expected gRPC codes are
`InvalidArgument`, `FailedPrecondition`, `Unavailable`, and `Internal`.
Health checks use the standard `grpc.health.v1.Health` service.

Go import path:

```go
import kvv1 "github.com/kirillidk/distributed-kv-storage/api/gen/go/kv/v1"
```

## Development

Go and [Buf](https://buf.build/docs/cli/) are required. Generator versions are
pinned in `buf.gen.yaml`.

```bash
make format
make generate
make ci
make breaking
```

Generated code is committed to the repository and must not be edited manually.

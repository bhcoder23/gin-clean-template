# AGENTS.md

This file applies to the entire repository.

## Project Identity

- This repository is the `bhcoder23`-maintained Gin variant of the original `evrone/go-clean-template`.
- Keep the repository positioned as a clean-architecture Go service template with multiple transports.
- This is a scaffold/template first. Demo business code exists to show correct boundaries, not to define product requirements.
- Preserve route behavior and public API shape unless the user explicitly asks for breaking changes.
- Default local development is HTTP-first. gRPC, AMQP RPC, NATS RPC, tracing, and outbox are optional adapters/examples.

## Architecture Rules

- Keep framework-specific code inside adapters only.
- `Gin` belongs only in `internal/transport/restapi` and `pkg/httpserver`.
- Do not introduce `gin.Context`, HTTP DTOs, or transport concerns into `internal/domain`, `internal/usecase`, or `internal/infra`.
- Business rules belong in `internal/usecase` and `internal/domain`, not in handlers or middleware.
- Repository and integration implementations should stay in `internal/infra/...`.
- Application-facing contracts should stay in `internal/usecase/contracts.go`.
- Keep `internal/usecase/contracts.go` grouped in this order:
  - inbound ports consumed by transports (`User`, `Task`, `Notification`)
  - outbound ports consumed by use case implementations and implemented by infrastructure (`UserRepo`, `TaskRepo`, `NotificationRepo`, `Transactor`, `RepoProvider`, `OutboxStore`)
  - port data used by those contracts (`TaskFilter`, `NotificationFilter`, `OutboxEvent`)
- Do not add framework types, persistence rows, or transport DTOs to `internal/usecase/contracts.go`.
- `internal/usecase` may depend on ports from `internal/usecase/contracts.go`; it must not import infrastructure packages.
- `internal/domain` should stay framework-free and persistence-free. Domain errors can live next to the model they describe.
- `pkg/...` is for reusable technical components. Do not put application business rules there.

## Scaffold Design Principles

- Prefer boring, explicit template code over clever abstractions.
- Avoid overfitting the scaffold to the demo task/user/notification business flow.
- When a demo feature teaches an architectural boundary, implement that boundary correctly. Examples: transaction boundaries, stable errors, request ID propagation, logging, and outbox publishing.
- Keep optional capabilities easy to remove. A derived project should be able to delete unused transports without rewriting use cases.
- Do not add a new layer or helper package unless it removes real duplication or protects a clear boundary.

## REST Conventions

- Follow the existing REST layout:
  - top-level router in `internal/transport/restapi/router.go`
  - versioned routes in `internal/transport/restapi/v1`
  - version route dependency wiring in `internal/transport/restapi/v1/router.go`
  - version route group registration in `internal/transport/restapi/v1/routes.go`
  - request DTOs in `internal/transport/restapi/v1/request`
  - response DTOs in `internal/transport/restapi/v1/response`
  - transport-only helpers in `internal/transport/restapi/middleware`
- Prefer Gin idioms consistently:
  - `ShouldBindJSON`
  - `c.Request.Context()`
  - `c.Set` / `c.Get`
  - `c.JSON` / `c.AbortWithStatusJSON`
- Keep handlers thin and explicit. Avoid adding extra abstraction layers unless they clearly reduce duplication.

## Transport Rules

- REST, gRPC, AMQP RPC, and NATS RPC are adapters over the same use cases.
- Transport routers should group dependencies in a `RouterDeps` struct instead of growing long constructor signatures.
- Keep transport validation at the adapter boundary, but keep core invariants in `internal/usecase` or `internal/domain` so every transport is safe.
- Transport-specific request and response DTOs must stay inside the corresponding `internal/transport/...` tree.
- Stable client-facing errors should be mapped through `internal/apperror`.
- Do not reintroduce transport-specific error helper packages like `errmap`, `errlog`, or `rpcerror`.
- RPC protocols should preserve stable error codes/messages and request ID metadata.

## Errors, Logging, and Request IDs

- Use `internal/apperror` for stable error code mapping across REST, gRPC, AMQP RPC, and NATS RPC.
- Expected application errors should log at warning level; unexpected errors should log at error level.
- Use `pkg/logger` for application logging. Do not use standard library logging in runtime packages.
- Default logging should be stdout-friendly for containers.
- Preserve request/correlation ID propagation:
  - HTTP uses `X-Request-ID`.
  - gRPC, AMQP RPC, and NATS RPC use metadata/header propagation.
  - REST error responses should include `request_id`.

## Transactions and Outbox

- Multi-repository use cases should use the `Transactor` port from `internal/usecase/contracts.go`.
- Transaction-scoped stores are exposed through `StoreProvider`.
- Business code that needs transactional event publishing should write through the `OutboxStore` port exposed by `StoreProvider.Outbox()`.
- Use cases must not import `internal/infra/outbox` directly.
- The outbox package is a production-shaped scaffold example: it supports transactional insert, pending claim, stale publishing lock recovery, retries, bounded publish timeout, and failure tracking.
- Core NATS publish + flush is not a durable broker acknowledgment. Keep documentation clear that durable event guarantees require replacing the publisher with JetStream, Kafka, RabbitMQ confirms, or another durable mechanism.

## Observability

- Tracing is optional and disabled by default.
- Keep OpenTelemetry setup in `pkg/observability` or adapter wiring. Do not put tracing SDK calls in domain or usecase code.
- Metrics and health/readiness endpoints should stay transport/runtime concerns.

## Generated Files

- Do not hand-edit generated sqlc outputs under `internal/infra/persistence/sqlc/*.go`.
- Do not hand-edit generated protobuf outputs under `docs/proto/v1/*.pb.go` and `docs/proto/v1/*_grpc.pb.go`.
- Do not hand-edit generated Swagger outputs unless the user explicitly asks:
  - `docs/docs.go`
  - `docs/swagger.json`
  - `docs/swagger.yaml`
- If persistence SQL or schema changes, regenerate sqlc code with:
  - `make sqlc`
- If REST annotations change, regenerate Swagger with:
  - `go tool swag init --parseDependency -g internal/transport/restapi/router.go`
- If proto definitions change, regenerate code with:
  - `make proto-v1`

## Formatting and Verification

- Format Go code with:
  - `go tool gofumpt -l -w .`
  - `go tool gci write . --skip-generated -s standard -s default`
- Preferred verification for normal changes:
  - `go test ./internal/... ./pkg/... ./config`
- Preferred lint:
  - `go tool golangci-lint run`
- If transport/runtime behavior changes and Docker is available, also run:
  - `make compose-up-integration-test`
- The repository expects Go toolchain auto-selection. If local Go is older than `go.mod`, run commands with `GOTOOLCHAIN=auto` or use `make` targets.

## Tests and Mocks

- When interface contracts in `internal/usecase/contracts.go` change, regenerate mocks in `internal/usecase` via:
  - `make mock`
- Keep tests close to the layer they verify.
- Prefer focused tests first, then broader verification.
- For transaction or outbox behavior, include tests that prove the intended boundary is used.
- For transport behavior changes, prefer at least one integration path when practical.

## Documentation

- Keep `README.md` and `README_CN.md` aligned when project positioning or framework choices change.
- Keep documentation accurate to the current implementation, especially around:
  - REST framework
  - commands in `Makefile`
  - generated artifacts
  - enabled-by-default transports
  - request ID, errors, tracing, transaction, and outbox behavior
- Do not keep tracked planning scratch files in `docs/plans`.

# Gin Clean Template

[🇨🇳 中文](README_CN.md)

General-purpose Clean Architecture template for Go backends, maintained by `bhcoder23`.

[![License](https://img.shields.io/badge/License-MIT-success)](LICENSE)
[![Maintainer](https://img.shields.io/badge/Maintainer-bhcoder23-1f6feb)](https://github.com/bhcoder23)

[![Web Framework](https://img.shields.io/badge/Gin-Web%20Framework-blue)](https://github.com/gin-gonic/gin)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Persistence](https://img.shields.io/badge/sqlc-Type--safe%20SQL-blue)](https://sqlc.dev/)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/prometheus/client_golang)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## Overview

This repository is the Gin-based backend scaffold maintained by `bhcoder23`.

The template focuses on code organization that stays useful after a project leaves the demo stage:

- keep domain and use case code independent from HTTP, gRPC, brokers, and database drivers
- keep transport adapters thin and replaceable
- make transaction, error, request ID, logging, tracing, and outbox boundaries explicit
- keep optional adapters easy to remove from derived projects

Inspired by the original MIT-licensed project:
- [evrone/go-clean-template](https://github.com/evrone/go-clean-template)

This template is one application process with multiple transport adapters:

- AMQP RPC (based on RabbitMQ as [transport](https://github.com/rabbitmq/amqp091-go)
  and [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- NATS RPC (based on NATS as [transport](https://github.com/nats-io/nats.go)
  and [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- gRPC ([gRPC](https://grpc.io/) framework based on protobuf)
- REST API ([Gin](https://github.com/gin-gonic/gin) framework)

The default local developer path starts with HTTP only. The other transports stay available as optional adapters so derived projects can opt in without carrying dependencies they do not need.

The template includes three domains to demonstrate multi-service architecture.
They are sample domains for the scaffold, not required product boundaries:

- **User Authentication** — registration, login, JWT-based authorization
- **Task Management** — CRUD operations with status transitions (todo, in_progress, done)
- **Notification Feed** — task activity notifications with read tracking

The demo domains can be exposed through all four transports (REST, gRPC, AMQP RPC, NATS RPC), but derived projects are expected to keep only the adapters they need.

## Content

- [Start here](#start-here)
- [Demo flow](#demo-flow)
- [Domains](#domains)
- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Adding a module](docs/add-new-module.md)
- [Dependency Injection](#dependency-injection)
- [Clean Architecture](#clean-architecture)

## Start here

Use the HTTP-first path first. It keeps the template easy to trim while still exercising the main scaffold:

```sh
# Start PostgreSQL for the HTTP-first local path
make compose-up

# Run migrations and start the enabled transports
make run
```

To inspect every demo adapter in one process, start the optional brokers first, then use `make run-all-transports`.

```sh
make compose-up-adapters
make run-all-transports
```

Once the app is running, the fastest way to understand the scaffold is to walk one complete REST flow end to end.

## Demo flow

Register a user:

```sh
curl -s http://127.0.0.1:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"johndoe","email":"john@example.com","password":"secret123"}'
```

Log in and capture the JWT:

```sh
TOKEN=$(
  curl -s http://127.0.0.1:8080/v1/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"email":"john@example.com","password":"secret123"}' | jq -r '.token'
)
```

Read the authenticated profile:

```sh
curl -s http://127.0.0.1:8080/v1/user/profile \
  -H "Authorization: Bearer $TOKEN"
```

Create a task:

```sh
curl -s http://127.0.0.1:8080/v1/tasks \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Ship the scaffold","description":"Exercise the happy path"}'
```

List tasks:

```sh
curl -s 'http://127.0.0.1:8080/v1/tasks?limit=10&offset=0' \
  -H "Authorization: Bearer $TOKEN"
```

List unread notifications generated by the task flow:

```sh
curl -s 'http://127.0.0.1:8080/v1/notifications?unread_only=true&limit=10&offset=0' \
  -H "Authorization: Bearer $TOKEN"
```

## Domains

The template includes three fully implemented domains, each demonstrated across the available transport adapters.

### User Authentication

Registration, login, and JWT-based authorization.

| Operation   | REST                     | gRPC                     |
|-------------|--------------------------|--------------------------|
| Register    | `POST /v1/auth/register` | `AuthService/Register`   |
| Login       | `POST /v1/auth/login`    | `AuthService/Login`      |
| Get profile | `GET /v1/user/profile`   | `AuthService/GetProfile` |

- Passwords hashed with bcrypt
- JWT tokens with configurable expiry
- Auth middleware on all transports

### Task Management

CRUD operations with a status state machine.

| Operation  | REST                         | gRPC                         |
|------------|------------------------------|------------------------------|
| Create     | `POST /v1/tasks`             | `TaskService/CreateTask`     |
| List       | `GET /v1/tasks`              | `TaskService/ListTasks`      |
| Get        | `GET /v1/tasks/:id`          | `TaskService/GetTask`        |
| Update     | `PUT /v1/tasks/:id`          | `TaskService/UpdateTask`     |
| Transition | `PATCH /v1/tasks/:id/status` | `TaskService/TransitionTask` |
| Delete     | `DELETE /v1/tasks/:id`       | `TaskService/DeleteTask`     |

- Status transitions: `todo` → `in_progress` → `done` (and `in_progress` → `todo`)
- Pagination with `limit`/`offset` and optional status filter
- Tasks scoped to the authenticated user

### Notification Feed

Task activity notifications persisted in PostgreSQL and exposed through every transport.

| Operation  | REST                              | gRPC                                     |
|------------|-----------------------------------|------------------------------------------|
| List       | `GET /v1/notifications`           | `NotificationService/ListNotifications`  |
| Mark read  | `PATCH /v1/notifications/:id/read`| `NotificationService/MarkNotificationRead` |

- Notifications are generated when tasks are created or moved through the status flow
- Unread filtering with `unread_only=true`
- Read tracking with `read_at`

## Quick start

### Local development

Docker is optional. `.env.example` starts HTTP only; gRPC, RabbitMQ RPC, and NATS RPC are opt-in. The Docker Compose demo stack sets those flags explicitly when it needs the full adapter set.

```sh
# PostgreSQL for the default HTTP-first path
make compose-up
# Run app with migrations
make run
```

To force all demo transports on regardless of your current `.env`, use:

```sh
make compose-up-adapters
make run-all-transports
```

### Integration tests

```sh
# DB, app + migrations, integration tests.
# These tests require the integration build tag and are usually run by Jenkins or another pipeline runner.
make compose-up-integration-test
```

### Full docker stack with reverse proxy

```sh
make compose-up-all
```

Check services in the full demo stack:

- AMQP RPC:
  - URL: `amqp://guest:guest@127.0.0.1:5672/`
  - Client Exchange: `rpc_client`
  - Server Exchange: `rpc_server`
- NATS RPC:
  - URL: `nats://guest:guest@127.0.0.1:4222/`
  - Server Exchange: `rpc_server`
- REST API:
  - http://app.lvh.me/healthz | http://127.0.0.1:8080/healthz
  - http://app.lvh.me/readyz | http://127.0.0.1:8080/readyz
  - http://app.lvh.me/metrics | http://127.0.0.1:8080/metrics
  - http://app.lvh.me/swagger | http://127.0.0.1:8080/swagger
- gRPC:
  - URL: `tcp://grpc.lvh.me:8081` | `tcp://127.0.0.1:8081`
  - [v1/auth.proto](docs/proto/v1/auth.proto)
  - [v1/task.proto](docs/proto/v1/task.proto)
  - [v1/notification.proto](docs/proto/v1/notification.proto)
- PostgreSQL:
  - `postgres://user:myAwEsOm3pa55@w0rd@127.0.0.1:5432/db`
- RabbitMQ:
  - http://rabbitmq.lvh.me | http://127.0.0.1:15672
  - Credentials: `guest` / `guest`
- NATS monitoring:
  - http://nats.lvh.me | http://127.0.0.1:8222/
  - Credentials: `guest` / `guest`

## Project structure

### `cmd/app/main.go`

Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.

### `config`

The twelve-factor app stores config in environment variables (often shortened to `env vars` or `env`). Env vars are easy
to change between deploys without changing any code; unlike config files, there is little chance of them being checked
into the code repo accidentally; and unlike custom config files, or other config mechanisms such as Java System
Properties, they are a language- and OS-agnostic standard.

Config: [config.go](config/config.go)

Example: [.env.example](.env.example)

Default local transport flags:

- `HTTP_ENABLED=true`
- `GRPC_ENABLED=false`
- `RMQ_ENABLED=false`
- `NATS_ENABLED=false`

`APP_ENV=production` adds guardrails: Swagger must be disabled and the sample JWT secret must be replaced.

Request correlation is part of the base scaffold:

- HTTP reads and writes `X-Request-ID`.
- gRPC, AMQP RPC, and NATS RPC use `x-request-id` metadata/header propagation.
- REST error responses include `request_id` so logs and client failures can be joined.

Optional tracing is available but disabled by default:

- `TRACE_ENABLED=false`
- `TRACE_EXPORTER=stdout`
- `TRACE_SERVICE_NAME=gin-clean-template`

The stdout exporter is intentionally concrete so the integration can be verified locally. Derived projects can replace the exporter with OTLP/collector wiring without changing handlers or use cases.

[docker-compose.yml](docker-compose.yml) uses `env` variables to configure services.

### `docs`

Swagger documentation. Auto-generated by [swag](https://github.com/swaggo/swag) library.
You don't need to correct anything by yourself.

[Add a New Business Module](docs/add-new-module.md) explains how to add product code without breaking scaffold boundaries.

#### `docs/proto`

Protobuf files. They are used to generate Go code for gRPC services.
The proto files are also used to generate documentation for gRPC services.
You don't need to correct anything by yourself.

### `integration-test`

Integration tests.
They are launched as a separate container, next to the application container.
They are excluded from normal `go test ./...` runs and require the `integration` build tag when run directly.

### `internal/app`

Runtime composition lives here. `Run` creates infrastructure, repositories, use cases, optional outbox relay, and enabled transport servers.

The application is one process that can start multiple adapters. Which servers are started is controlled by:

- `HTTP_ENABLED`
- `GRPC_ENABLED`
- `RMQ_ENABLED`
- `NATS_ENABLED`

Shutdown is coordinated through an application root context. If the wiring grows beyond simple constructor calls, a DI generator such as [wire](https://github.com/google/wire) can be introduced later, but the scaffold keeps explicit constructors by default.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
go run -tags migrate ./cmd/app
```

### `internal/transport`

Incoming adapter layer. The template includes 4 optional transports:

- AMQP RPC (based on RabbitMQ as transport)
- NATS RPC (based on NATS as transport)
- gRPC ([gRPC](https://grpc.io/) framework based on protobuf)
- REST API ([Gin](https://github.com/gin-gonic/gin) framework)

Server routers are written in the same style:

- Handlers are grouped by application area
- Version router dependencies are grouped in a `RouterDeps` struct instead of long function signatures
- Route groups are registered explicitly in the version package
- Business logic interfaces are injected into the router controller, which handlers call

#### `internal/transport/amqp_rpc`

AMQP request/reply adapter backed by RabbitMQ. Routes are registered in `amqp_rpc/v1/routes.go`; auth binding, request validation, error mapping, and request ID propagation stay inside this adapter.

#### `internal/transport/grpc`

gRPC adapter backed by generated protobuf code under `docs/proto/v1`. Stable application errors are mapped to gRPC status codes and carry the client-facing code as `google.rpc.ErrorInfo.reason`.

#### `internal/transport/nats_rpc`

NATS request/reply adapter. It follows the same controller, route, request, and response layout as AMQP RPC so optional transports remain easy to compare or remove.

#### `internal/transport/restapi`

Gin REST adapter. Top-level runtime endpoints such as `/healthz`, `/readyz`, `/metrics`, and `/swagger/*any` are registered in `internal/transport/restapi/router.go`; versioned API routes live under `internal/transport/restapi/v1`.

REST request DTOs live in `v1/request`, response DTOs live in `v1/response`, and route groups are registered in `v1/routes.go`. Swagger annotations live next to the handlers and are generated with [swag](https://github.com/swaggo/swag).

### `internal/domain`

Core domain models and the rules that belong to them.
This layer contains entities, enums, value objects, and domain errors that should stay independent from transport and storage concerns.

### `internal/usecase`

Application business logic.

- Use case implementations live in domain-oriented subpackages such as `internal/usecase/task`
- Application-facing contracts live in `internal/usecase/contracts.go`
- Use cases move `internal/domain` values across business boundaries

Use cases depend on contracts defined in `internal/usecase/contracts.go`.
Persistence implementations, transport adapters, and reusable technical packages are injected into use cases
(see [Dependency Injection](#dependency-injection)).

#### `internal/infra/persistence`

Persistence implementations for PostgreSQL-backed repositories used by the use case layer.

### `pkg`

Reusable technical components live here. They are not allowed to contain application business rules.

Current examples include:

- `pkg/httpserver` and `pkg/grpcserver` server wrappers
- `pkg/postgres` connection pool and transaction executor helpers
- `pkg/logger` zerolog adapter
- `pkg/requestid` request/correlation ID helpers
- `pkg/rabbitmq` and `pkg/nats` request/reply primitives
- `pkg/observability` optional OpenTelemetry setup

## Dependency Injection

The scaffold uses explicit constructor injection. The goal is not to hide dependencies; it is to make boundaries visible.

The normal direction is:

- `internal/app` creates concrete infrastructure and use case implementations
- `internal/usecase` defines the contracts it needs
- `internal/infra/...` implements outbound contracts such as repositories and outbox storage
- `internal/transport/...` consumes inbound use case contracts

For example, the task use case receives repository ports and an optional transaction port:

```go
task.New(taskRepo, notificationRepo, transactor)
```

Tests use generated mocks from `internal/usecase/contracts.go`:

```sh
make mock
```

If wiring grows too large, introduce [wire](https://github.com/google/wire) or another DI generator at the application boundary. Keep the contracts and dependency direction unchanged.

## Clean Architecture

### Current Rules

The repository follows a pragmatic ports-and-adapters layout:

- `internal/domain` contains framework-free domain models, domain errors, enums, and model-level rules
- `internal/usecase` contains application workflows and the contracts those workflows need
- `internal/infra` contains concrete outbound adapters, currently PostgreSQL persistence and outbox storage
- `internal/transport` contains inbound adapters, currently REST, gRPC, AMQP RPC, and NATS RPC
- `pkg` contains reusable technical components, not product business rules

The key dependency rule is simple: inner layers do not import outer layers.

Allowed examples:

- transport -> usecase contract -> usecase implementation
- usecase implementation -> usecase contract -> infra implementation
- infra implementation -> domain conversion

![Clean Architecture](docs/img/layers-1.png)

Do not put framework types, request DTOs, persistence rows, broker clients, or database handles into `internal/domain` or use case contracts.

### Boundary Example

For an HTTP request that needs database data, the flow is:

```
    HTTP > usecase
           usecase > persistence contract
           usecase < persistence contract
    HTTP < usecase
```

The symbols > and < show layer boundaries crossed through interfaces.

![Example](docs/img/example-http-db.png)

For a workflow that also publishes events, the flow should remain explicit:

```
    HTTP > usecase
           usecase > persistence contract
           usecase < persistence contract
           usecase > external integration contract
           usecase < external integration contract
           usecase > RPC
           usecase < RPC
           usecase > persistence contract
           usecase < persistence contract
    HTTP < usecase
```

![Example](docs/img/layers-2.png)

### Domain, DTO, and Persistence Rows

Domain values are the language of business code. Transports and persistence adapters convert at the edge:

- REST request/response structs stay under `internal/transport/restapi/v1/request` and `response`
- AMQP/NATS request/response structs stay under their transport trees
- gRPC protobuf messages stay under `docs/proto/v1`
- sqlc rows stay under `internal/infra/persistence/sqlc`
- conversion between rows and domain values is owned by the repository implementation

### Transactions and Persistence

For cross-repository writes, the persistence layer exposes a small transaction template instead of requiring an ORM:

- `persistence.NewRepositories(pg)` creates repositories backed by the normal pool.
- `persistence.NewTransactor(pg).WithinTx(ctx, fn)` creates repositories backed by one `pgx` transaction.
- Repositories depend on the minimal `postgres.Executor` interface, so the same repository can run on a pool or a transaction.
- Demo repository SQL lives in `internal/infra/persistence/sql/*.sql` and is generated with `sqlc` into `internal/infra/persistence/sqlc`.

This is intentionally a template extension point. Simple single-repository demo use cases can call repositories directly; flows that need atomic multi-table updates should opt into `WithinTx` without leaking `pgx.Tx` into `internal/usecase`. The task sample uses this boundary when writing the task and its notification. `sqlc` removes hand-written scan boilerplate while the repository still owns domain conversion and stable error mapping.

### Errors

REST errors use a stable envelope:

```json
{
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "task not found",
    "request_id": "..."
  }
}
```

The mapping is centralized in `internal/apperror`, following the same idea as Kratos' `code`/`reason` split: transport status codes remain protocol-level, while the string `code` is the stable client-facing reason. REST responses expose it as `error.code`; gRPC responses attach it as `google.rpc.ErrorInfo.reason`; AMQP RPC and NATS RPC use it as the RPC status code. Demo domain errors live next to their sample model files; REST, gRPC, AMQP, and NATS error mapping plus expected-error log classification are handled in `apperror` to avoid duplicated transport helper packages.

### Outbox

For event publishing, the scaffold includes a production-shaped transactional outbox example:

- migration-backed `outbox_events` table
- `outbox.Store` for transactional inserts, pending claims, and stale publishing-lock recovery
- `outbox.Relay` with retries, lock timeout, bounded publish timeout, and failure tracking
- `outbox.NATSPublisher` as the concrete default publisher binding with client-side flush

It is disabled by default through `OUTBOX_ENABLED=false`. When enabled with `OUTBOX_PUBLISHER=nats`, the relay publishes events to `OUTBOX_SUBJECT_PREFIX + "." + event_type`. Business use cases should write outbox rows through the `OutboxStore` port exposed by the transaction `RepoProvider`, inside the same `WithinTx` callback as their database changes, when they need DB + outbox consistency. Core NATS publish + flush confirms the client handed the message to the server connection; it is not a durable broker acknowledgment. Swap the publisher to JetStream, Kafka, RabbitMQ confirms, or another durable mechanism when the business event must survive broker-side failure.

## Similar projects

- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## Useful links

- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/)

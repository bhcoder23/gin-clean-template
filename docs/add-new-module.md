# Add a New Business Module

This guide shows where code should go when a derived project adds a new business module.

The scaffold is intentionally adapter-first at the edges and domain-first in the core:

- HTTP, gRPC, AMQP RPC, and NATS RPC code stay in `internal/transport/...`.
- Business rules stay in `internal/usecase/...` and `internal/domain/...`.
- Persistence and external integrations stay in `internal/infra/...`.
- Application-facing contracts stay in `internal/usecase/contracts.go`.
- Reusable technical helpers can live in `pkg/...`; application business code should not.

Use the existing task/user/notification demo only as an architectural example. Do not copy domain details unless your product actually needs them.

## Recommended Steps

### 1. Define the Domain

Add framework-free model and domain errors under `internal/domain`.

Example:

```text
internal/domain/project.go
```

Domain code must not import Gin, gRPC, SQL, Redis, NATS, RabbitMQ, or DTO packages.

Keep stable domain errors close to the model when they describe core business rules:

```go
var (
    ErrProjectNotFound = errors.New("project not found")
    ErrProjectNameRequired = errors.New("project name is required")
)
```

### 2. Add Use Case Contracts

Extend `internal/usecase/contracts.go` with the ports the use case needs.

Typical additions:

- a repository interface, such as `ProjectRepo`
- filters or command structs if they clarify signatures
- a use case interface, such as `Project`

Keep the contract application-facing. Do not mention infrastructure types such as `*gorm.DB`, `pgx.Tx`, `*sql.DB`, `*redis.Client`, or broker clients.

If a flow writes multiple repositories, use the existing `Transactor` port. If it writes business data and integration events atomically, write through `RepoProvider.Outbox()`.

### 3. Implement the Use Case

Create a package under `internal/usecase/<module>`.

Example:

```text
internal/usecase/project/project.go
```

Use cases should accept `context.Context`, validate core invariants, call ports from `internal/usecase/contracts.go`, and return domain errors for expected failures.

Do not pass `*gin.Context` into use cases. Transport adapters should convert request data into ordinary Go values before calling business logic.

### 4. Implement Persistence

Add repository code under `internal/infra/persistence`.

Example:

```text
internal/infra/persistence/project_postgres.go
```

The repository implements the repo interface from `internal/usecase/contracts.go`.

If the module needs transaction-scoped repositories, update the repository provider in `internal/infra/persistence` so `RepoProvider` can expose the repository inside `WithinTx`.

### 5. Add Migrations

Add explicit migrations under `migrations`.

Do not rely on automatic schema migration in application startup. The scaffold uses explicit migration files so deploys are auditable and repeatable.

### 6. Map Stable Errors

Register expected domain errors in `internal/apperror`.

Each stable application error should define:

- HTTP status
- gRPC code
- client-facing string code
- safe client message
- whether the error is expected

This keeps REST, gRPC, AMQP RPC, and NATS RPC consistent.

### 7. Add REST Adapter Code

Use the existing REST layout:

```text
internal/transport/restapi/v1/project.go
internal/transport/restapi/v1/project_routes.go
internal/transport/restapi/v1/request/project.go
internal/transport/restapi/v1/response/project.go
```

Handlers should stay thin:

- bind and validate transport DTOs
- read request-scoped metadata, such as user ID
- call the use case with `c.Request.Context()`
- map results to response DTOs
- delegate expected errors to the shared error response path

Do not put business rules in handlers.

### 8. Add Optional RPC Adapters

Only add gRPC, AMQP RPC, or NATS RPC adapters when the product needs them.

Keep request/response DTOs inside the corresponding transport tree:

```text
internal/transport/grpc/v1
internal/transport/amqp_rpc/v1
internal/transport/nats_rpc/v1
```

If you add proto definitions, update files under `docs/proto/v1` and run:

```sh
make proto-v1
```

If you add REST Swagger annotations, run:

```sh
go tool swag init --parseDependency -g internal/transport/restapi/router.go
```

### 9. Wire the Module in `internal/app`

Create repositories and use cases in `internal/app/app.go` or a nearby app wiring file when the module is enabled.

Prefer explicit constructor calls until the wiring becomes too large. Wire or another DI generator can be introduced later, but it should solve real wiring complexity rather than hide simple setup.

### 10. Test the Boundary

Add tests close to the layer being verified:

- domain tests for business invariants
- use case tests for transaction, outbox, and port behavior
- repository tests for persistence queries
- transport tests for request validation, error mapping, request ID, and status codes

When a contract in `internal/usecase/contracts.go` changes, regenerate mocks:

```sh
make mock
```

Run the normal verification:

```sh
go test ./internal/... ./pkg/... ./config
```

## Naming Checklist

For a module named `project`, a typical REST-first implementation looks like this:

```text
internal/domain/project.go
internal/usecase/project/project.go
internal/infra/persistence/project_postgres.go
internal/transport/restapi/v1/project.go
internal/transport/restapi/v1/project_routes.go
internal/transport/restapi/v1/request/project.go
internal/transport/restapi/v1/response/project.go
migrations/<timestamp>_create_projects.up.sql
migrations/<timestamp>_create_projects.down.sql
```

Update these shared files only when the module needs them:

```text
internal/usecase/contracts.go
internal/infra/persistence/stores.go
internal/apperror/error.go
internal/app/app.go
```

## What Not to Do

- Do not put application routes, controllers, DTOs, or use cases under `pkg`.
- Do not pass `gin.Context` into domain, use case, or infrastructure code.
- Do not import infrastructure packages from `internal/usecase`.
- Do not return raw database or broker errors directly to clients.
- Do not add a framework or generator just because a module was added.
- Do not enable optional transports unless the product needs them.

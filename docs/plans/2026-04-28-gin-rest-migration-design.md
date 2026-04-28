# Gin REST Migration Design

**Goal:** Replace the current Fiber-based REST adapter with a Gin-based REST adapter in `gin-clean-template`, while keeping the clean architecture boundaries, existing HTTP API surface, and non-REST transports unchanged.

**Scope:**
- Replace `Fiber` with `Gin` only in `pkg/httpserver` and `internal/controller/restapi`
- Remove Fiber-specific dependencies from `go.mod`
- Keep request/response DTOs, use cases, repositories, gRPC, NATS RPC, RabbitMQ RPC, and app wiring intact
- Preserve existing routes, status codes, auth behavior, Swagger docs, metrics endpoint, and graceful shutdown behavior

**Architecture decisions:**
- `Gin` stays confined to the REST adapter boundary so that `entity`, `usecase`, `repo`, and `app` remain framework-agnostic.
- `pkg/httpserver` becomes a thin wrapper around `gin.Engine` and `net/http.Server`.
- REST handlers continue to translate HTTP concerns into use case calls; no business logic moves into handlers.
- Shared REST helper style is standardized around `ShouldBindJSON`, `c.JSON`, `c.Set/Get`, and `c.Request.Context()`.

**Style rules for the migration:**
- Keep framework types out of non-REST layers.
- Keep route groups and handler structure parallel to the current Fiber layout.
- Keep error mapping explicit inside handlers.
- Prefer small helpers for repeated HTTP concerns, but avoid new abstractions unless they reduce obvious duplication.
- Use standard library HTTP primitives where Gin naturally interoperates with them.

**Dependency choices:**
- `github.com/gin-gonic/gin` at latest stable version
- `github.com/swaggo/gin-swagger` and `github.com/swaggo/files` for Swagger UI
- `github.com/prometheus/client_golang/prometheus/promhttp` for the `/metrics` endpoint via `gin.WrapH`

**Testing strategy:**
- Update the existing auth middleware test first so it targets Gin behavior and fails before production changes.
- Run focused REST/package tests during migration.
- Run broader project tests after the REST adapter is fully migrated.

**Non-goals:**
- No route redesign
- No use case or repository refactor
- No gRPC or MQ transport changes
- No extra platform features beyond the REST framework swap

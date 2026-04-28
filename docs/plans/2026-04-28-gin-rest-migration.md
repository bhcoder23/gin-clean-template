# Gin REST Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace the Fiber REST adapter in `gin-clean-template` with a Gin REST adapter while preserving API behavior and clean architecture boundaries.

**Architecture:** The migration is constrained to `pkg/httpserver` and `internal/controller/restapi`. Gin becomes the only REST framework dependency, while `internal/entity`, `internal/usecase`, `internal/repo`, and the non-REST transports remain unchanged. Metrics, Swagger, auth, logging, recovery, and graceful shutdown stay available through Gin-compatible integrations.

**Tech Stack:** Go 1.26, Gin v1.12.0, gin-swagger, promhttp, validator, zerolog-based app logger.

---

### Task 1: Migrate auth middleware test to Gin

**Files:**
- Modify: `internal/controller/restapi/middleware/auth_test.go`

**Step 1: Write the failing test**
- Replace Fiber app setup with Gin engine setup.
- Keep the same auth cases: missing header, bad scheme, invalid token, valid token.

**Step 2: Run test to verify it fails**
Run: `go test ./internal/controller/restapi/middleware -run TestAuthMiddleware -v`
Expected: FAIL because production middleware still exposes Fiber types.

**Step 3: Write minimal implementation**
- None in this task.

**Step 4: Run test to verify it still fails for the expected reason**
Run: `go test ./internal/controller/restapi/middleware -run TestAuthMiddleware -v`
Expected: FAIL with type/compile errors pointing at Fiber-only middleware.

### Task 2: Replace HTTP server wrapper with Gin

**Files:**
- Modify: `pkg/httpserver/server.go`
- Modify: `pkg/httpserver/options.go`

**Step 1: Write the failing test**
- Reuse the middleware test red state from Task 1.

**Step 2: Implement server wrapper**
- Replace `fiber.App` with `gin.Engine`
- Replace framework-specific shutdown/start logic with `net/http.Server`
- Keep options and graceful shutdown semantics aligned with current behavior

**Step 3: Run focused tests**
Run: `go test ./internal/controller/restapi/middleware -run TestAuthMiddleware -v`
Expected: still FAIL until middleware is migrated.

### Task 3: Migrate REST middleware to Gin

**Files:**
- Modify: `internal/controller/restapi/middleware/auth.go`
- Modify: `internal/controller/restapi/middleware/logger.go`
- Modify: `internal/controller/restapi/middleware/recovery.go`

**Step 1: Implement Gin auth middleware**
- Use `gin.HandlerFunc`
- Store `userID` through `c.Set`
- Abort with JSON on auth failures

**Step 2: Implement Gin logger middleware**
- Log after `c.Next()` with IP, method, URL, status, and response size

**Step 3: Implement Gin recovery middleware**
- Use `gin.CustomRecovery`
- Log panic with stack trace and return `500` JSON payload

**Step 4: Run focused tests**
Run: `go test ./internal/controller/restapi/middleware -run TestAuthMiddleware -v`
Expected: PASS

### Task 4: Migrate routers and handlers to Gin

**Files:**
- Modify: `internal/controller/restapi/router.go`
- Modify: `internal/controller/restapi/v1/router.go`
- Modify: `internal/controller/restapi/v1/error.go`
- Modify: `internal/controller/restapi/v1/task.go`
- Modify: `internal/controller/restapi/v1/user.go`
- Modify: `internal/controller/restapi/v1/translation.go`

**Step 1: Replace route registration with Gin groups**
- Keep route paths unchanged
- Keep protected/public split unchanged

**Step 2: Replace handler context usage**
- `BodyParser` -> `ShouldBindJSON`
- `Locals` -> `Get`
- `UserContext()` -> `Request.Context()`
- `Status(...).JSON(...)` -> `JSON(...)`

**Step 3: Wire Swagger and metrics for Gin**
- Swagger: `gin-swagger`
- Metrics: `promhttp.Handler()` wrapped with Gin

**Step 4: Run focused package tests**
Run: `go test ./internal/controller/restapi/... ./pkg/httpserver -v`
Expected: PASS

### Task 5: Remove Fiber dependencies and verify formatting

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Possibly modify generated Swagger imports if needed

**Step 1: Update dependencies**
- Add latest Gin and Gin Swagger deps
- Remove Fiber and Fiber-specific middleware deps
- Run `go mod tidy`

**Step 2: Format code**
Run: `gofumpt -l -w . && gci write . --skip-generated -s standard -s default`
Expected: no remaining style drift

**Step 3: Run project verification**
Run: `go test ./internal/... ./pkg/... -v`
Expected: PASS

### Task 6: Regenerate REST docs if adapter imports changed

**Files:**
- Modify if regenerated: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`

**Step 1: Regenerate Swagger docs**
Run: `swag init --parseDependency -g internal/controller/restapi/router.go`
Expected: generated docs remain valid and compile

**Step 2: Run final verification**
Run: `go test ./internal/... ./pkg/... -v`
Expected: PASS

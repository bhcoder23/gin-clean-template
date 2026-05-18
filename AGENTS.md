# AGENTS.md

This file applies to the entire repository.

## Project Identity

- This repository is the `bhcoder23`-maintained Gin variant of the original `evrone/go-clean-template`.
- Keep the repository positioned as a clean-architecture Go service template with multiple transports.
- Preserve route behavior and public API shape unless the user explicitly asks for breaking changes.

## Architecture Rules

- Keep framework-specific code inside adapters only.
- `Gin` belongs only in `internal/transport/restapi` and `pkg/httpserver`.
- Do not introduce `gin.Context`, HTTP DTOs, or transport concerns into `internal/domain`, `internal/usecase`, or `internal/infra`.
- Business rules belong in `internal/usecase` and `internal/domain`, not in handlers or middleware.
- Repository and integration implementations should stay in `internal/infra/...`.
- Application-facing contracts should stay in `internal/usecase/contracts.go`.

## REST Conventions

- Follow the existing REST layout:
  - top-level router in `internal/transport/restapi/router.go`
  - versioned routes in `internal/transport/restapi/v1`
  - request DTOs in `internal/transport/restapi/v1/request`
  - response DTOs in `internal/transport/restapi/v1/response`
  - transport-only helpers in `internal/transport/restapi/middleware`
- Prefer Gin idioms consistently:
  - `ShouldBindJSON`
  - `c.Request.Context()`
  - `c.Set` / `c.Get`
  - `c.JSON` / `c.AbortWithStatusJSON`
- Keep handlers thin and explicit. Avoid adding extra abstraction layers unless they clearly reduce duplication.

## Generated Files

- Do not hand-edit generated protobuf outputs under `docs/proto/v1/*.pb.go` and `docs/proto/v1/*_grpc.pb.go`.
- Do not hand-edit generated Swagger outputs unless the user explicitly asks:
  - `docs/docs.go`
  - `docs/swagger.json`
  - `docs/swagger.yaml`
- If REST annotations change, regenerate Swagger with:
  - `go tool swag init --parseDependency -g internal/transport/restapi/router.go`
- If proto definitions change, regenerate code with:
  - `make proto-v1`

## Formatting and Verification

- Format Go code with:
  - `go tool gofumpt -l -w .`
  - `go tool gci write . --skip-generated -s standard -s default`
- Preferred verification for normal changes:
  - `go test ./internal/... ./pkg/...`
- If transport/runtime behavior changes and Docker is available, also run:
  - `make compose-up-integration-test`

## Tests and Mocks

- When interface contracts in `internal/usecase/contracts.go` change, regenerate mocks in `internal/usecase` via:
  - `make mock`
- Keep tests close to the layer they verify.
- Prefer focused tests first, then broader verification.

## Documentation

- Keep `README.md` and `README_CN.md` aligned when project positioning or framework choices change.
- Keep documentation accurate to the current implementation, especially around:
  - REST framework
  - commands in `Makefile`
  - generated artifacts

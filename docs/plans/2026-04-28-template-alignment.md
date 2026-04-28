# Template Alignment Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Finish the ownership, naming, tooling, and generated-doc alignment for `gin-clean-template` so it can serve as a polished general-purpose backend scaffold.

**Architecture:** The change is intentionally shallow: identity and developer-experience cleanup only. Existing clean-architecture layering, sample domains, and transport adapters stay intact. Generated Swagger artifacts are refreshed after source annotations and module imports are updated.

**Tech Stack:** Go 1.26, Gin v1.12.0, swaggo, Prometheus client, zerolog, Make.

---

### Task 1: Align repository identity

**Files:**
- Modify: `go.mod`
- Modify: repo-wide Go imports referencing the old module path
- Modify: `.env.example`
- Modify: `docker-compose.yml`
- Modify: `Makefile`

**Step 1: Update module path**
- Change the module path to `github.com/bhcoder23/gin-clean-template`.

**Step 2: Update internal imports**
- Replace all `github.com/evrone/go-clean-template/...` imports with the new module path.

**Step 3: Update template naming**
- Replace remaining `go-clean-template` app-name and volume-name defaults with `gin-clean-template`.

### Task 2: Align tooling commands

**Files:**
- Modify: `Makefile`

**Step 1: Switch to `go tool`-based commands**
- Use `go tool swag`, `go tool gofumpt`, `go tool gci`, `go tool golangci-lint`, and `go tool govulncheck` where applicable.

**Step 2: Keep existing workflows intact**
- Preserve target names and expected developer ergonomics.

### Task 3: Refresh docs and generated artifacts

**Files:**
- Modify: `README.md`
- Modify: `README_CN.md`
- Modify: `README_RU.md`
- Modify: `internal/controller/restapi/router.go`
- Regenerate: `docs/docs.go`
- Regenerate: `docs/swagger.json`
- Regenerate: `docs/swagger.yaml`

**Step 1: Update positioning text**
- Keep the project framed as a general-purpose backend template maintained by `bhcoder23`.

**Step 2: Update Swagger metadata**
- Replace old template title strings with Gin-aligned naming.

**Step 3: Regenerate generated docs**
- Run Swagger generation after source annotations are updated.

### Task 4: Cleanup and verify

**Files:**
- Modify: `docker-compose.yml`
- Modify: `docker-compose-integration-test.yml`

**Step 1: Remove obsolete Compose `version` keys**
- Keep service definitions unchanged otherwise.

**Step 2: Run formatting and verification**
- Run `go mod tidy`
- Run `go tool gofumpt -l -w .`
- Run `go tool gci write . --skip-generated -s standard -s default`
- Run `go test ./internal/... ./pkg/...`
- Run `go test -race ./internal/... ./pkg/...`
- Run `go tool golangci-lint run`

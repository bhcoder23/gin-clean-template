# Template Alignment Design

**Goal:** Align `gin-clean-template` with its current GitHub identity and make it a cleaner general-purpose backend scaffold without changing the existing architecture, transports, or sample domains.

**Scope:**
- Rename the Go module path to `github.com/bhcoder23/gin-clean-template`
- Replace remaining template identity strings such as `go-clean-template`
- Align local tooling commands with the Go `tool` block
- Remove obsolete Docker Compose `version` fields
- Refresh generated Swagger metadata and related docs wording

**Recommended approach:**
- Keep the current project structure, sample domains, and transport coverage unchanged
- Treat this as a consistency pass, not a framework redesign
- Apply identity and tooling cleanup repo-wide, then regenerate generated artifacts and re-run verification

**Non-goals:**
- No change to clean architecture boundaries
- No change to REST/gRPC/RabbitMQ/NATS behavior
- No new production features, infra, or runtime dependencies

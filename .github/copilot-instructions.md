# Copilot Instructions for FastQuest Backend

## Build, test, and lint

```bash
# Run the API locally (expects DB/Gemini env vars)
go run .

# Build all packages
go build ./...

# Run all tests/packages (current baseline has no *_test.go files yet)
go test ./...

# Run a single test (when present)
go test ./internal/question -run TestCreateQuestion -v
```

There is no repository-specific linter config (`golangci-lint`, `Makefile`, or task runner) checked in. Use Go defaults when needed:

```bash
gofmt -w .
go vet ./...
```

## High-level architecture

- Entry point is `main.go`: it initializes the database (`internal/platform/database.InitDB`), initializes Gemini (`internal/ai.InitGemini`), then starts the HTTP server from `NewServer()`.
- Route registration is centralized in `main.go` (`registerPaths`), mapping directly to domain handlers under `internal/{question,answer,questionset,source,exam,ai}`.
- Each domain follows a layered package shape (`handler.go`, `service.go`, `repository.go`, `dto.go`):
  - **handler**: HTTP parsing/validation/response
  - **repository**: GORM queries and persistence helpers
  - **service**: internal orchestration/validation methods used across modules (especially by AI flows)
  - **dto**: request payload types
- Persistence uses a shared global GORM connection from `internal/platform/database` (`database.GetDB()`), backed by PostgreSQL.
- Shared data contracts live in `pkg/models`; response shaping is done via `ToResponse()` methods and include scopes like `ApplyQuestionIncludes` / `ApplyQuestionSetIncludes`.
- API contracts are maintained in two places:
  - Swagger artifacts under `docs/` (`swagger.yaml`, `swagger.json`, generated `docs.go`)
  - Human-oriented route documentation in `ROTAS.md`

## Key codebase conventions

- Database table names are singular and explicitly mapped with `TableName()` in models (`question`, `answer`, `question_set`, etc.). Keep model/table naming aligned with existing migrations.
- `include` query params are comma-separated and mapped to GORM preloads through `Apply*Includes` scope functions in `pkg/models` (for example: `include=answers,user,subject,source`).
- List endpoints consistently use pagination query params `page` and `perPage` (with max `perPage` 100) and return a `data + pagination` JSON envelope.
- Question creation (`POST /questions`) intentionally accepts either:
  - a single question object, or
  - a non-empty array of question objects (batch create).
- Write operations that span multiple tables are expected to use DB transactions (see `CreateQuestionSet` and `CreateSource` handlers).
- Environment loading convention:
  - local/dev: `.env` is loaded when `RAILWAY_ENVIRONMENT` is not set
  - deployed environments: rely on injected env vars
  - expected keys are documented in `.env.example` (DB, Goose, Gemini).

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Kiosk-based student check-in web app. Go backend + static HTML/JS frontend served from the same binary. Postgres for storage. Spec lives in `.kiro/specs/student-check-in-system/`.

## Where to look

Six canonical docs. Anything else `*.md` in the repo root is stale session notes — ignore it.

- `README.md` — public overview, install/run, project layout
- `CLAUDE.md` — this file; rules for AI agents working in this repo
- `CHANGELOG.md` — dated, human-curated history of notable changes
- `ARCHITECTURE_SIMPLIFIED.md` — why there's no service layer (Handler → Repository → DB)
- `LOGGING.md` — how to use `internal/logger`, rotation, caller-depth gotchas
- `TESTING.md` — unit test conventions, `./run_tests.sh`, `./test_api.sh` flow

For change history use `git log` — it's more reliable than any narrative doc.

## Common commands

```bash
# Bring up Postgres (docker-compose also defines an unused mssql service)
docker-compose up -d postgres

# Run the server (loads .env, runs migrations, listens on :8090)
go run cmd/server/main.go

# Unit tests (validation + handlers, with coverage at the end)
./run_tests.sh

# Or run tests directly
go test ./internal/...
go test -v -run TestValidateStudent ./internal/validation/...
go test ./internal/... -coverprofile=coverage.out && go tool cover -html=coverage.out

# End-to-end API smoke test (requires server running on :8090)
./test_api.sh
```

The server hard-codes port `:8090` in `cmd/server/main.go` — `SERVER_PORT` in `.env` is not currently read. DB env vars (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`) are required and validated at startup in `internal/db/postgres.go`.

## Architecture

Two-layer design (intentionally — see `ARCHITECTURE_SIMPLIFIED.md`): **Handler → Repository → DB**. There is no service layer. Handlers own HTTP I/O *and* business logic (validation, uniqueness checks, orchestration). Don't reintroduce a service layer for a single-caller helper.

Request flow:
1. `cmd/server/main.go` initializes logger → DB → migrations → handlers → router, then `http.ListenAndServe`.
2. `internal/router/router.go` (chi) wires all routes under `/api/*` and serves `web/static/` at `/*`. Custom `requestLogger` middleware logs every request with method/path/status/duration; `/favicon.ico` is logged at DEBUG, Chrome devtools probes are silenced.
3. Handlers (`internal/handlers/`) decode JSON, call `validation.Sanitize*` then `validation.Validate*`, call repository methods, and respond via the `respondWithJSON` / `respondWithError` helpers in `student_handler.go`.
4. Repositories (`internal/repository/`) are interfaces with a `Postgres*Repository` implementation. Use `context.Context` from the request and `QueryRowContext`/`QueryContext`. Duplicate-key errors are detected by string-matching `"duplicate key value"` / `"unique constraint"`.

### Domain entities

- **Student** (`models/student.go`) — has optional `ClassID *int` and a legacy `ClassInfo string`. Both can coexist; `ClassInfo` is kept for backward compat.
- **Class** (`models/class.go`) — has optional `StudentID *int` (the class leader, who is also a student). Migration `008_rename_leader_to_student.sql` renamed this from `leader_id`; if you see `leader_id` references in old docs, the current name is `student_id`.
- **Event**, **Checkin** — straightforward.
- **User** (`models/user.go`, `handlers/user.go`) is legacy from the original scaffold and still wired at `/api/user`, `/api/users`, `/api/user/{id}`. Don't extend it for new check-in features; use Student/Event/Checkin.

### Migrations

`internal/migration/migration.go` is a homegrown forward-only runner: it scans `migrations/*.sql` lexicographically, tracks applied files in a `schema_migrations` table by filename, and runs anything new. There is no down-migration support and `golang-migrate/migrate` (in `go.sum`) is **not** actually used. New migrations: prefix with the next zero-padded number (`009_…sql`); the whole file runs as a single `db.Exec`, so don't include `\` psql meta-commands.

### Logging

`internal/logger/logger.go` wraps the stdlib `log` with lumberjack rotation to `logs/server.log` (10MB, 10 backups, 30 days, gzipped). Use `logger.Info/Warn/Error/Debug(msg, key, value, ...)` with key/value pairs — `formatLog` walks args in pairs. `logger.Request(...)` is reserved for the router middleware. The logger uses `runtime.Caller(3)` to auto-tag file:line and function name; if you wrap these calls in another helper, the caller depth will be wrong.

### Frontend

Plain HTML/JS in `web/static/` — no build step. Pages live in `pages/`, scripts in `js/`, served directly by the chi file server. The check-in flow expects the API at the same origin.

## Testing notes

Unit tests live alongside code (`*_test.go`) in `internal/validation/` and `internal/handlers/`. Repository and DB layers are not unit-tested — exercise them via `./test_api.sh` against a running server. Validation tests are table-driven; follow that pattern for new validators.

## Conventions

- Markdown docs: only the six listed under **Where to look** are canonical. Do not create `*_SUMMARY.md`, `*_FIX.md`, `*_COMPLETE.md`, `IMPROVEMENTS_*.md`, `FIXES_ROUND*.md`, `TEST_RESULTS.md`, or other session-artifact docs. Findings from a session belong in the commit message, the PR description, or this conversation — not a new top-level doc.
- `.env` is gitignored; `.env.example` and `.env.production` are tracked templates.
- Backup files (`main.go.bak`, `main.go.bak2`, `data.bak/`, `mssql-data/`, `postgres-data/`) are leftovers — leave them unless cleanup is requested.

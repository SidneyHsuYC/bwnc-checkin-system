# bwnc-checkin-system

Kiosk-based student check-in web app. Single Go binary that serves a JSON API at `/api/*` and the static frontend (`web/static/`) from `/`. Postgres for storage.

## Quick start

```bash
# 1. Bring up Postgres (docker-compose also defines an unused mssql service)
docker-compose up -d postgres

# 2. Copy and edit env
cp .env.example .env

# 3. Run the server (loads .env, runs migrations, listens on :8090)
go run cmd/server/main.go
```

Open <http://localhost:8090> for the kiosk UI.

## Browsing the database

`docker-compose.yml` includes an [Adminer](https://www.adminer.org/) container for ad-hoc SQL and table browsing.

```bash
docker-compose up -d adminer
```

Then open <http://localhost:8080> and log in with:

| Field    | Value                |
|----------|----------------------|
| System   | PostgreSQL           |
| Server   | `checkin-postgres`   |
| Username | `checkin_user`       |
| Password | `checkin_pass`       |
| Database | `checkin_db`         |

(`Server` is the Postgres container name on the docker network, not `localhost` — Adminer is connecting from inside the docker network.) Default credentials match `docker-compose.yml`; if you've overridden them, use yours.

## Configuration

Set in `.env` (see `.env.example` for the full template):

| Var           | Required | Notes                                       |
|---------------|----------|---------------------------------------------|
| `DB_HOST`     | yes      |                                             |
| `DB_PORT`     | yes      |                                             |
| `DB_USER`     | yes      |                                             |
| `DB_PASSWORD` | yes      |                                             |
| `DB_NAME`     | yes      |                                             |
| `DB_SSLMODE`  | yes      | `disable` for local                         |
| `SERVER_PORT` | no       | Currently ignored — port is hard-coded to `:8090` in `cmd/server/main.go` |

## API

All JSON. All under `/api/*`.

| Method | Path                       | Description                                    |
|--------|----------------------------|------------------------------------------------|
| GET    | `/health`                  | Liveness + DB connectivity                     |
| POST   | `/api/students`            | Create student                                 |
| GET    | `/api/students/search?q=`  | Search by name/email                           |
| POST   | `/api/classes`             | Create class                                   |
| GET    | `/api/classes`             | List classes                                   |
| GET    | `/api/classes/search?q=`   | Search classes                                 |
| GET    | `/api/classes/{id}`        | Get class (with leader info if set)            |
| PUT    | `/api/classes/{id}`        | Update class                                   |
| DELETE | `/api/classes/{id}`        | Delete class                                   |
| POST   | `/api/events`              | Create event                                   |
| GET    | `/api/events`              | List events                                    |
| GET    | `/api/events/recent`       | Recent events for the kiosk picker             |
| POST   | `/api/checkins`            | Record a check-in (rejects duplicates)         |
| GET    | `/api/checkins`            | List check-ins                                 |
| POST   | `/api/user`, `GET /api/users`, `GET /api/user/{id}` | Legacy user scaffold — don't extend |

The legacy `/api/user(s)` endpoints are kept for backward compatibility from the original scaffold; new check-in features should use Student/Event/Checkin.

## Project layout

```
cmd/server/main.go         # entry point: env → logger → DB → migrations → router
internal/
├── db/                    # Postgres connection + retry/pool
├── migration/             # forward-only migration runner (homegrown)
├── models/                # Student, Class, Event, Checkin, User
├── validation/            # Sanitize* + Validate* helpers (unit-tested)
├── repository/            # interface + Postgres impl per entity
├── handlers/              # HTTP layer + business logic (no service layer)
├── router/                # chi router + requestLogger middleware
└── logger/                # stdlib log + lumberjack rotation
migrations/                # 001_…sql, 002_…sql — applied lexicographically
web/static/                # plain HTML/JS/CSS, no build step
```

## Testing

```bash
./run_tests.sh           # Go unit tests (validation + handlers)
./test_api.sh            # API smoke test — requires server running on :8090
```

See `TESTING.md` for conventions and what each layer covers.

## Further reading

- `CLAUDE.md` — rules and orientation for AI agents working in this repo
- `ARCHITECTURE_SIMPLIFIED.md` — why there's no service layer
- `LOGGING.md` — logger usage and gotchas
- `CHANGELOG.md` — notable changes

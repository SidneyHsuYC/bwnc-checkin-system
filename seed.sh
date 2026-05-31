#!/usr/bin/env bash
#
# seed.sh — load example data (seed.sql) into the check-in database.
#
# Run it once after the server has created the tables on first start:
#
#     docker-compose up -d postgres
#     go run cmd/server/main.go          # creates the tables, then Ctrl-C
#     ./seed.sh
#
# It's safe to run again — seed.sql is idempotent (it never deletes or resets
# data). Connection details come from .env (falling back to the same defaults
# as .env.example / docker-compose.yml). If a local `psql` is installed it's
# used; otherwise the script runs psql inside the docker-compose Postgres
# container, so you don't need any Postgres client on your laptop.

set -euo pipefail

cd "$(dirname "$0")"

SEED_FILE="seed.sql"
CONTAINER="checkin-postgres"   # matches docker-compose.yml

# Load DB_* from .env if present (without clobbering anything already exported).
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-checkin_user}"
DB_PASSWORD="${DB_PASSWORD:-checkin_pass}"
DB_NAME="${DB_NAME:-checkin_db}"

if [ ! -f "$SEED_FILE" ]; then
  echo "❌ $SEED_FILE not found (run this from the repo root)." >&2
  exit 1
fi

echo "Seeding example data into '$DB_NAME'..."

if command -v psql >/dev/null 2>&1; then
  echo "  using local psql -> ${DB_USER}@${DB_HOST}:${DB_PORT}"
  PGPASSWORD="$DB_PASSWORD" psql \
    -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -v ON_ERROR_STOP=1 -f "$SEED_FILE"
elif docker exec "$CONTAINER" true >/dev/null 2>&1; then
  echo "  using docker container '$CONTAINER'"
  docker exec -i "$CONTAINER" psql \
    -U "$DB_USER" -d "$DB_NAME" \
    -v ON_ERROR_STOP=1 < "$SEED_FILE"
else
  echo "❌ No way to reach Postgres." >&2
  echo "   Install the psql client, or start the database with:" >&2
  echo "       docker-compose up -d postgres" >&2
  exit 1
fi

echo "✅ Done. Open http://localhost:8090 and try a check-in."

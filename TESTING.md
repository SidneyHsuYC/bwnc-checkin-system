# Testing

Two layers, run separately:

| Layer        | Runner          | Requires                | What it covers                          |
|--------------|-----------------|-------------------------|------------------------------------------|
| Unit         | `./run_tests.sh` (or `go test ./internal/...`) | nothing            | `internal/validation/`, `internal/handlers/` |
| API smoke    | `./test_api.sh` | server running on :8090 | full HTTP → DB round trip                |

Repository and DB layers have no unit tests — they're exercised through `test_api.sh`.

## Unit tests

```bash
./run_tests.sh                                        # all packages, with coverage summary
go test ./internal/validation/...                     # one package
go test -v -run TestValidateStudent ./internal/...    # one test
go test ./internal/... -coverprofile=coverage.out \
  && go tool cover -html=coverage.out                 # HTML coverage report
```

Tests live next to the code (`*_test.go`) in:

- `internal/validation/` — `student_validation_test.go`, `event_validation_test.go`, `checkin_validation_test.go`
- `internal/handlers/` — `student_handler_test.go`

### Conventions

- **Table-driven**, always. New validators should follow the same shape as existing ones:

  ```go
  func TestValidateX(t *testing.T) {
      tests := []struct {
          name    string
          input   *models.X
          wantErr bool
      }{
          {name: "valid", input: &models.X{...}, wantErr: false},
          {name: "missing field", input: &models.X{...}, wantErr: true},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              err := ValidateX(tt.input)
              if (err != nil) != tt.wantErr {
                  t.Errorf("ValidateX() error = %v, wantErr %v", err, tt.wantErr)
              }
          })
      }
  }
  ```

- Test names describe the scenario (`"missing first name"`, `"email too long"`), not the function under test.
- Cover the error branches, not just the happy path. The validators have a small finite set of error cases — exhaust them.
- Don't introduce DB/network dependencies into unit tests. If you need the database, write an API test instead.

## API smoke test

```bash
# in one terminal
docker-compose up -d postgres
go run cmd/server/main.go

# in another
./test_api.sh
```

Prereqs: server up on `:8090`, `jq` installed (`brew install jq`).

The script walks the public surface — class CRUD, student registration with and without a class, event creation, check-in, duplicate-email rejection. It prints `✓` / `✗` per step and exits non-zero on the first failure.

If the script fails:

- `curl http://localhost:8090/health` to confirm the server is reachable
- check that `:8090` isn't held by something else: `lsof -ti:8090`
- if the DB has stale state, `docker-compose down -v` will wipe the volume and the next server start will re-run migrations from scratch

## Manual UI testing

There's no automated frontend test. Smoke checks worth running after touching `web/static/`:

- Student registration: email autocomplete suggests `@gmail.com`/`@yahoo.com`/etc; class autocomplete pulls from the `/api/classes` endpoint; submitting with an empty class field still succeeds.
- Check-in flow: "Don't find your name? Register now!" link works; after a successful check-in, the primary action is "Go Back to Home Page" with "Continue to Check-in" secondary.
- Duplicate check-in for the same (student, event) should be rejected by the backend.

## Coverage targets

Validation is the only package where 100% is realistic and expected — the logic is small and fully deterministic. Handlers sit lower because the DB-touching paths aren't in scope for unit tests. Don't chase a coverage number; chase covering each error branch in `internal/validation/` and each non-DB code path in `internal/handlers/`.

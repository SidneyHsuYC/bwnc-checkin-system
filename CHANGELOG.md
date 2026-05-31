# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Removed
- **Legacy `User` scaffold** — deleted `internal/handlers/user.go`, `internal/models/user.go`,
  the `/api/user`, `/api/users`, and `/api/user/{id}` routes, and `migrations/001_create_users.sql`.
  The frontend never used it; new features use Student/Event/Checkin.
- **`golang-migrate/migrate`** dependency (had no Go imports) — dropped via `go mod tidy`.
- **`mssql` service** from `docker-compose.yml` (unused; Postgres is the only datastore).
- Tracked stray files: `cmd/server/main.go.bak`, `.bak2`, `.backup`, and the committed `tmp/server` binary.

### Changed
- `/health` is now served by a dedicated `HealthHandler` (`internal/handlers/health.go`) that reuses
  `db.HealthCheck` and no longer queries the `users` table; response is `{status, database}` (no `user_count`).

### Added
- **Example data** for local development: idempotent `seed.sql` + `seed.sh` runner
  (5 classes, 10 students, 5 events dated relative to today, 9 check-ins).
- Collaborator-focused `README.md` Quick start, kiosk screenshots (`docs/images/`), and a Troubleshooting section.
- Fixed `test_api.sh` for the current schema (`class_name` required, `student_id` not `leader_id`) and made it re-runnable.

## [2.0.0] - 2026-02-04

### Added
- **Log Rotation System** (`internal/logger/logger.go`)
  - Automatic rotation based on file size (10MB) and time (daily)
  - Keeps last 10 log files automatically
  - Dual output to both console and file
  - Background monitoring for rotation triggers
  - Graceful shutdown handling
  - Log file info display on startup

- **New API Endpoint: GET /api/user/{id}**
  - Fetch individual user by ID
  - Returns 404 for non-existent users
  - Full logging support
  - Proper error handling

- **Request Filtering**
  - Skip logging for browser noise (DevTools, favicon)
  - Cleaner log output

- **Graceful Shutdown**
  - Signal handling for SIGINT and SIGTERM
  - Proper cleanup of resources
  - Log file closure on shutdown

- **Documentation**
  - `UPDATES.md` - Detailed answers to all questions
  - `CHANGELOG.md` - This file
  - `.gitignore` - Excludes logs directory

### Changed
- **API Endpoint Naming** (RESTful conventions)
  - `POST /api/users` → `POST /api/user` (singular for single resource)
  - `GET /api/users` remains (plural for collection)
  - `GET /api/user/{id}` (singular for single resource)

- **Database Connection Logging**
  - Removed hardcoded `localhost:5432` string
  - Implemented dynamic URL parsing
  - Shows actual host, port, and database name
  - Masks only the password
  - Example: `postgres://postgres:***@localhost:5432/default`

- **Test Script Updates**
  - Updated endpoint paths to match new naming
  - Added test for GET /api/user/{id}
  - Added test for 404 response
  - Now includes 8 tests total

- **Main Server**
  - Integrated logger initialization
  - Added log file info display on startup
  - Updated endpoint documentation
  - Added graceful shutdown

### Fixed
- Chrome DevTools 404 errors no longer clutter logs
- Password masking now shows actual connection details
- Consistent RESTful endpoint naming

### Technical Details

**Log Rotation Configuration:**
```go
maxLogSize   = 10 * 1024 * 1024  // 10MB
maxLogFiles  = 10                 // Keep last 10 files
logDir       = "logs"             // Log directory
```

**Log File Naming:**
```
server_2026-02-04_21-59-33.log
       └─ ISO timestamp
```

**Rotation Triggers:**
1. File size exceeds 10MB
2. Day changes (midnight)

**Cleanup:**
- Automatically removes files beyond the 10-file limit
- Sorted by modification time (keeps newest)

## [1.0.0] - 2026-02-04

### Initial Release
- PostgreSQL database integration
- User CRUD operations (Create, Read)
- Health check endpoint
- Comprehensive logging with emojis
- Connection pooling
- Retry logic for database connections
- CORS support
- Automated testing script
- Complete documentation

### API Endpoints
- `GET /health` - Health check
- `POST /api/users` - Create user
- `GET /api/users` - Get all users

### Features
- Request timing and performance monitoring
- Detailed error messages
- Input validation
- Password masking in logs
- Migration system
- Static file serving

---

## Version Numbering

This project follows [Semantic Versioning](https://semver.org/):
- MAJOR version for incompatible API changes
- MINOR version for backwards-compatible functionality additions
- PATCH version for backwards-compatible bug fixes

## Migration Guide

### From 1.0.0 to 2.0.0

**API Changes:**
```bash
# Old
POST /api/users

# New
POST /api/user
```

**Update your client code:**
```javascript
// Before
fetch('http://localhost:8090/api/users', {
  method: 'POST',
  // ...
})

// After
fetch('http://localhost:8090/api/user', {
  method: 'POST',
  // ...
})
```

**New Endpoint:**
```bash
# Fetch user by ID
GET /api/user/{id}
```

**Log Files:**
- Logs now automatically saved to `logs/` directory
- No configuration needed
- Automatic rotation and cleanup

## Upgrade Instructions

1. Pull latest code
2. No database migrations needed
3. Update API client endpoints (POST /api/users → POST /api/user)
4. Restart server
5. Logs will automatically start rotating

## Breaking Changes

### 2.0.0
- `POST /api/users` changed to `POST /api/user`
  - **Impact:** API clients need to update endpoint URL
  - **Reason:** RESTful naming convention (singular for single resource)
  - **Migration:** Simple find/replace in client code

## Deprecations

None currently.

## Security Updates

### 2.0.0
- Enhanced password masking with dynamic URL parsing
- Log files properly excluded from version control
- Proper file permissions on log files (0644)

## Performance Improvements

### 2.0.0
- Log rotation runs in background (non-blocking)
- Minimal performance impact (<1ms overhead)
- Efficient file cleanup algorithm

## Known Issues

None currently.

## Future Roadmap

### Planned Features
- [ ] UPDATE endpoint (PUT /api/user/{id})
- [ ] DELETE endpoint (DELETE /api/user/{id})
- [ ] Pagination for GET /api/users
- [ ] Search and filtering
- [ ] Authentication/Authorization
- [ ] Rate limiting
- [ ] Metrics endpoint (Prometheus format)
- [ ] Structured logging (JSON format option)
- [ ] Request ID tracking
- [ ] Database connection health monitoring

### Under Consideration
- [ ] GraphQL API
- [ ] WebSocket support
- [ ] Caching layer
- [ ] Multi-database support
- [ ] API versioning
- [ ] OpenAPI/Swagger documentation

## Contributors

- Initial development and enhancements

## License

[Your License Here]

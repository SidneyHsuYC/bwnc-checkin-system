# Logging

`internal/logger/logger.go` wraps the stdlib `log` package and writes through a `lumberjack.Logger` for size-based rotation, plus a duplicate sink to stdout.

## Configuration

Defined as constants in `logger.go`:

| Setting       | Value             | Notes                                  |
|---------------|-------------------|----------------------------------------|
| `logDir`      | `logs/`           | Created on `Init()` if missing         |
| `logFile`     | `logs/server.log` | Active log file                        |
| `maxLogSize`  | 10 MB             | Rotation trigger                       |
| `maxBackups`  | 10                | Older backups deleted automatically    |
| `maxAge`      | 30 days           | Older backups deleted automatically    |
| `compressOld` | `true`            | Rotated files become `.gz`             |

Rotated files are named like `server-2026-02-04T10-30-15.log` (or `.log.gz` once compressed). To change any of these, edit the constants — there is no runtime config.

## API

```go
logger.Init()           // call once from main, after parsing env
defer logger.Close()    // optional; just emits a shutdown line

logger.Info(msg, key, value, ...)
logger.Warn(msg, key, value, ...)
logger.Error(msg, key, value, ...)
logger.Debug(msg, key, value, ...)
```

Variadic args are walked as key/value pairs by `formatLog`. An odd number of args drops the trailing key — pass them in pairs.

`logger.Request(level, method, path, status, duration)` is reserved for the chi `requestLogger` middleware in `internal/router/router.go`. Don't call it from handlers.

## Output format

```
student_handler.go:142	[INFO]	[CreateStudent]	student created id=17 email=alice@example.com
```

Tab-delimited fields:
1. `file:line` — caller location
2. `[LEVEL]`
3. `[FunctionName]` — captured automatically (see below)
4. message + space-separated `key=value` pairs

The line is prefixed with the date/time set by `log.SetFlags(log.Ldate | log.Ltime)`.

## Automatic function name capture

`getCaller()` uses `runtime.Caller(3)` to walk three frames up the stack:

```
runtime.Caller(3) → caller of Info/Error/Warn/Debug
                 → formatLog        (frame 1)
                 → Info/Error/etc.  (frame 2)
                 → user code        (frame 3, what we want)
```

This means **don't wrap `logger.Info` etc. in a helper** — if you do, frame 3 will be your helper, not the real caller, and every log line will report the wrong function name. If you need a wrapper, either build it on top of `formatLog` directly with a tuned `runtime.Caller(N)`, or accept the wrong tag.

`logger.Request` does its own caller capture and skips the indirection level, so it's safe to call from middleware.

## Operations

```bash
# tail current log
tail -f logs/server.log

# search compressed backups
zgrep ERROR logs/server-*.log.gz

# disk usage
du -sh logs/
```

If logs aren't appearing: check that `logger.Init()` was called and `logs/` is writable. If rotation isn't happening: confirm the file actually crossed `maxLogSize` (10 MB) — rotation is size-triggered, not time-triggered.

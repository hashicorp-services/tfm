# Structured Logging with log/slog

## Why slog Over logrus

`log/slog` is Go's stdlib structured logging package (since Go 1.21). It replaces third-party loggers like logrus (maintenance mode since 2020) and zap.

**Benefits of slog:**

- Zero external dependencies
- Structured key-value pairs by design
- Pluggable handlers (`TextHandler`, `JSONHandler`, custom)
- Runtime-mutable log levels via `slog.LevelVar`
- `AddSource: true` replaces manual `runtime.Caller` hacks
- Direct use as dependency — `*slog.Logger` IS the interface, no wrapper needed

**Anti-pattern: Custom Logger interfaces wrapping slog.** Don't create `type Logger interface { Debug(msg string, args ...any) }` — just use `*slog.Logger` directly. It already is a clean, well-designed interface. Custom wrappers block slog's handler ecosystem and add indirection for no benefit.

## Setup

### Basic Logger with LevelVar

```go
func buildLogger(level string) (*slog.Logger, *slog.LevelVar) {
    levelVar := &slog.LevelVar{}

    // Map level string to slog level
    switch strings.ToLower(level) {
    case "trace", "debug":
        levelVar.Set(slog.LevelDebug)
    case "", "info":
        levelVar.Set(slog.LevelInfo)
    case "warning", "warn":
        levelVar.Set(slog.LevelWarn)
    case "error", "fatal", "panic", "critical":
        levelVar.Set(slog.LevelError)
    default:
        levelVar.Set(slog.LevelInfo)
    }

    handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
        Level:     levelVar,
    })

    return slog.New(handler), levelVar
}
```

**Key points:**

- `slog.LevelVar` enables runtime level changes without rebuilding the logger
- `AddSource: true` automatically adds `source=file.go:42` — no `runtime.Caller` needed
- Store `levelVar` alongside logger for commands that need `ApplyLogLevel`

### Runtime Level Changes

```go
func ApplyLogLevel(level string, lv *slog.LevelVar) error {
    if level == "" {
        return nil
    }

    switch strings.ToLower(level) {
    case "trace", "debug":
        lv.Set(slog.LevelDebug)
    case "info":
        lv.Set(slog.LevelInfo)
    case "warning", "warn":
        lv.Set(slog.LevelWarn)
    case "error", "fatal", "panic", "critical":
        lv.Set(slog.LevelError)
    default:
        return fmt.Errorf("invalid log level %q", level)
    }
    return nil
}
```

**Backward compatibility note:** slog's `Level.UnmarshalText` only recognizes `DEBUG`, `INFO`, `WARN`, `ERROR`. If migrating from logrus, add a pre-mapping for logrus level names like `trace`, `warning`, `fatal`, `panic`, `notice`, `critical`.

## Structured Logging Patterns

### Use Structured Attributes, Not fmt.Sprintf

```go
// BAD: Buries structured data in formatted string
logger.Info(fmt.Sprintf("Scheduler started with %d jobs", jobCount))

// GOOD: Structured attributes enable machine parsing and filtering
logger.Info("Scheduler started", "jobCount", jobCount)

// BAD: Error details lost in string formatting
logger.Error(fmt.Sprintf("Job %s failed: %v", name, err))

// GOOD: Each field is independently queryable
logger.Error("Job failed", "job", name, "error", err)
```

### Key Naming Conventions

```go
// Use camelCase for attribute keys (Go convention)
logger.Info("Request completed",
    "method", r.Method,
    "path", r.URL.Path,
    "statusCode", resp.StatusCode,
    "duration", time.Since(start),
)

// Group related attributes
logger.Info("Job completed",
    slog.Group("job",
        slog.String("name", job.GetName()),
        slog.String("type", "exec"),
    ),
    slog.Group("execution",
        slog.Duration("duration", d),
        slog.Bool("failed", false),
    ),
)
```

### Passing Loggers Through Structs

```go
// Use *slog.Logger directly in struct fields — it IS the interface
type Scheduler struct {
    Logger   *slog.Logger
    LevelVar *slog.LevelVar // Only if runtime level changes needed
    // ...
}

type Context struct {
    Logger    *slog.Logger
    Execution *Execution
    Job       Job
}

// Create child loggers with additional context
func (s *Scheduler) runJob(job Job) {
    jobLogger := s.Logger.With("job", job.GetName())
    jobLogger.Info("Starting job")
    // ...
    jobLogger.Info("Job completed", "duration", elapsed)
}
```

## Middleware Logging

```go
// Logging middleware using slog
func WithLogging(logger *slog.Logger) Middleware {
    return func(next Job) Job {
        return JobFunc(func(ctx context.Context) error {
            start := time.Now()
            logger.Info("Starting job", "job", next.GetName())

            err := next.Run(ctx)

            attrs := []any{
                "job", next.GetName(),
                "duration", time.Since(start),
            }
            if err != nil {
                logger.Error("Job failed", append(attrs, "error", err)...)
            } else {
                logger.Info("Job completed", attrs...)
            }
            return err
        })
    }
}
```

## Web Handler Logging

```go
type Handler struct {
    scheduler *Scheduler
    logger    *slog.Logger
}

func (h *Handler) TriggerJob(w http.ResponseWriter, r *http.Request) {
    name := chi.URLParam(r, "name")
    reqLogger := h.logger.With("handler", "TriggerJob", "jobName", name)

    job, err := h.scheduler.GetJob(name)
    if err != nil {
        reqLogger.Warn("Job not found")
        http.Error(w, "Job not found", http.StatusNotFound)
        return
    }

    go func() {
        if err := job.Run(context.Background()); err != nil {
            reqLogger.Error("Manual job execution failed", "error", err)
        }
    }()

    w.WriteHeader(http.StatusAccepted)
}
```

## Testing with slog

### Discard Logger for Tests

```go
// Simple: discard all logs
logger := slog.New(slog.NewTextHandler(io.Discard, nil))

// With specific level (only errors logged)
logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
    Level: slog.LevelError,
}))
```

### Capturing Logs in Tests

For tests that need to assert on log output, implement a custom `slog.Handler`:

```go
type TestHandler struct {
    mu      sync.Mutex
    records []slog.Record
}

func (h *TestHandler) Enabled(_ context.Context, _ slog.Level) bool {
    return true
}

func (h *TestHandler) Handle(_ context.Context, r slog.Record) error {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.records = append(h.records, r.Clone())
    return nil
}

func (h *TestHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *TestHandler) WithGroup(_ string) slog.Handler      { return h }

// Query helpers
func (h *TestHandler) HasMessage(msg string) bool {
    h.mu.Lock()
    defer h.mu.Unlock()
    for _, r := range h.records {
        if strings.Contains(r.Message, msg) {
            return true
        }
    }
    return false
}

func (h *TestHandler) HasAttr(key, value string) bool {
    h.mu.Lock()
    defer h.mu.Unlock()
    for _, r := range h.records {
        r.Attrs(func(a slog.Attr) bool {
            if a.Key == key && strings.Contains(a.Value.String(), value) {
                return false // found, stop iteration
            }
            return true
        })
    }
    return false
}
```

**Important:** The test handler captures `r.Message` separately from attributes. If your test assertions check for values that were moved from `fmt.Sprintf` to structured attributes during a migration, you need to update assertions to check attributes, not message strings.

```go
// Usage in tests
func TestRetryLogging(t *testing.T) {
    handler := &TestHandler{}
    logger := slog.New(handler)

    retrier := NewRetrier(logger)
    retrier.Execute(failingFunc)

    // Check message text
    assert.True(t, handler.HasMessage("Job failed, retrying"))
    // Check structured attributes
    assert.True(t, handler.HasAttr("attempt", "1"))
    assert.True(t, handler.HasAttr("maxRetries", "3"))
}
```

## Migration from logrus

### Level Mapping

| logrus | slog | Notes |
|--------|------|-------|
| `Trace` | `Debug` | slog has no Trace; use Debug |
| `Debug` | `Debug` | Direct mapping |
| `Info` | `Info` | Direct mapping |
| `Warn` / `Warning` | `Warn` | Direct mapping |
| `Error` | `Error` | Direct mapping |
| `Fatal` | `Error` + `os.Exit(1)` | slog has no Fatal; log then exit |
| `Panic` | `Error` + `panic()` | slog has no Panic; log then panic |

### Callsite Conversion

```go
// logrus printf-style
logger.Debugf("loaded config from %s", path)
logger.WithField("job", name).WithError(err).Error("execution failed")
logger.WithFields(logrus.Fields{"job": name, "attempt": n}).Warn("retrying")

// slog structured style
logger.Debug("loaded config", "file", path)
logger.Error("execution failed", "job", name, "error", err)
logger.Warn("retrying", "job", name, "attempt", n)
```

### Migration Checklist

1. Replace `Logger` interface/field types with `*slog.Logger` throughout
2. Replace `logrus.New()` with `slog.New(slog.NewTextHandler(...))`
3. Convert `logger.Debugf("msg %s", x)` to `logger.Debug("msg", "key", x)`
4. Replace `logrus.Fields{...}` with inline key-value pairs
5. Replace `WithError(err)` with `"error", err` attribute
6. Replace `WithField("k", v)` with `logger.With("k", v)` for persistent context
7. Delete custom Logger interfaces — `*slog.Logger` IS the interface
8. Delete logrus adapter/wrapper code
9. Update test loggers (see Testing section above)
10. Run `go mod tidy` to remove logrus from go.mod
11. Add logrus to depguard deny list in `.golangci.yml`
12. **Update CI workflows** — remove references to deleted logging packages

### CI Gotcha

When deleting a logging package, check CI workflow files for references:

```yaml
# BAD: References deleted package — CI will fail with [setup failed]
go test -race ./core/... ./config/... ./logging/...

# GOOD: Removed deleted package
go test -race ./core/... ./config/...

# BAD: Matrix includes deleted package
package: [cli, core, config, logging, middlewares, web]

# GOOD: Removed from matrix
package: [cli, core, config, middlewares, web]
```

Always grep CI workflows after deleting any package:
`grep -r "package-name" .github/workflows/`

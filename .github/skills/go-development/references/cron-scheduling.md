# Cron Scheduling with go-cron

[`github.com/netresearch/go-cron`](https://github.com/netresearch/go-cron) is a maintained fork of `robfig/cron` — the most popular cron library for Go — with bug fixes, runtime schedule updates, per-entry context, resilience middleware, and modern toolchain support.

## Installation

```bash
go get github.com/netresearch/go-cron
```

```go
import cron "github.com/netresearch/go-cron"
```

Drop-in replacement for `robfig/cron/v3` — just change the import path.

## Basic Usage

```go
c := cron.New()

c.AddFunc("0 9 * * *", func() {
    fmt.Println("Every day at 9am")
})

c.AddFunc("@every 5m", func() {
    fmt.Println("Every 5 minutes")
})

c.Start()
defer c.Stop()
```

## Named Jobs and Lookup

Assign names and tags for O(1) lookup, update, and removal:

```go
id, _ := c.AddFunc("0 9 * * *", dailyReport,
    cron.WithName("daily-report"),
    cron.WithTags("reports", "daily"),
)

// Lookup by name (O(1))
entry := c.EntryByName("daily-report")

// Filter by tag
entries := c.EntriesByTag("reports")

// Remove by name
c.RemoveByName("daily-report")
```

## Runtime Updates

Update schedules and jobs atomically without remove+re-add:

```go
// Update schedule only (preserves job, options, and context)
c.UpdateScheduleByName("daily-report", cron.Every(5*time.Minute))

// Update both schedule and job atomically (cancels old entry context)
c.UpdateEntryJobByName("daily-report", "30 10 * * *", newJob)

// Create-or-update in one call
id, err := c.UpsertJob("0 9 * * *", myJob, cron.WithName("my-job"))
```

### Graceful Job Replacement

For long-running jobs, wait for the current execution to finish before replacing:

```go
c.WaitForJobByName("my-job")  // Block until current execution finishes
c.UpsertJob(newSpec, newJob, cron.WithName("my-job"))
```

Check if a job is currently running:

```go
if c.IsJobRunningByName("my-job") {
    log.Println("Job is still running, will wait")
    c.WaitForJobByName("my-job")
}
```

## Per-Entry Context

Each entry gets its own `context.Context` derived from the Cron's base context. The context is automatically canceled when the entry is removed or its job is replaced.

```go
c.AddJob("@every 1m", cron.FuncJobWithContext(func(ctx context.Context) {
    select {
    case <-ctx.Done():
        return // Entry removed or job replaced
    case <-time.After(10 * time.Second):
        // Work completed
    }
}))
```

### Context Hierarchy

```
caller's context
  └─ cron context (canceled by Stop())
       └─ entry context (canceled by Remove/UpdateEntry/UpsertJob)
```

`cron.New(cron.WithContext(parentCtx))` derives a child context. `Stop()` cancels the child, not the caller's context.

## Job Wrappers (Middleware)

### Concurrency Wrappers

These implement `JobWithContext` and propagate context to inner jobs:

```go
// Apply to all jobs via Cron options
c := cron.New(cron.WithChain(
    cron.Recover(logger),              // Catch panics
    cron.SkipIfStillRunning(logger),   // Skip if previous still running
    cron.DelayIfStillRunning(logger),  // Queue until previous finishes
    cron.Timeout(30*time.Second, nil), // Abandon after duration
    cron.TimeoutWithContext(30*time.Second, nil), // Cancel context after duration
    cron.Jitter(5*time.Second),        // Random delay
))

// Apply to specific job
job := cron.NewChain(
    cron.Recover(logger),
    cron.DelayIfStillRunning(logger),
).Then(myJob)
```

### Resilience Wrappers

These return `FuncJob` and do NOT forward context:

```go
// Retry on panic with exponential backoff
retryJob := cron.RetryWithBackoff(myJob, cron.RetryConfig{
    MaxRetries:   3,
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:     30 * time.Second,
    Multiplier:   2.0,
})

// Retry on error return (job must implement ErrorJob)
retryJob := cron.RetryOnError(myErrorJob, cron.RetryOnErrorConfig{
    MaxRetries: 3,
    Delay:      time.Second,
})

// Circuit breaker — stop after consecutive failures
cbJob := cron.CircuitBreaker(myJob, cron.CircuitBreakerConfig{
    Threshold:    5,
    ResetTimeout: time.Minute,
})
```

### ErrorJob Interface

For retry-on-error, implement `ErrorJob`:

```go
type myJob struct{}

func (j *myJob) Run() {}
func (j *myJob) RunWithError() error {
    // Return error to trigger retry
    return doWork()
}
```

Or use the convenience wrapper:

```go
cron.FuncErrorJob(func() error {
    return doWork()
})
```

## Observability

Monitor cron operations with hooks:

```go
c := cron.New(cron.WithObservability(cron.ObservabilityHooks{
    OnJobStart: func(id cron.EntryID, name string, scheduled time.Time) {
        jobsStarted.WithLabelValues(name).Inc()
    },
    OnJobComplete: func(id cron.EntryID, name string, dur time.Duration, recovered any) {
        jobDuration.WithLabelValues(name).Observe(dur.Seconds())
        if recovered != nil {
            jobPanics.WithLabelValues(name).Inc()
        }
    },
}))
```

## Validation

Validate cron expressions before scheduling:

```go
// Quick validation
if err := cron.ValidateSpec("0 9 * * MON-FRI"); err != nil {
    log.Fatal(err)
}

// Instance-level (uses configured parser)
c := cron.New(cron.WithSeconds())
if err := c.ValidateSpec("0 30 * * * *"); err != nil {
    log.Fatal(err)
}

// Detailed analysis
result := cron.AnalyzeSpec("0 9 * * MON-FRI")
fmt.Println("Next run:", result.NextRun)
fmt.Println("Fields:", result.Fields)
```

## Missed Job Catch-Up

Handle jobs missed during downtime:

```go
lastRun := loadFromDatabase("daily-report")

c.AddFunc("0 9 * * *", dailyReport,
    cron.WithPrev(lastRun),
    cron.WithMissedPolicy(cron.MissedRunOnce),
    cron.WithMissedGracePeriod(2*time.Hour),
)
```

Policies: `MissedSkip` (default), `MissedRunOnce`, `MissedRunAll`.

## Graceful Shutdown

```go
// Block until all running jobs finish
c.StopAndWait()

// With timeout
if !c.StopWithTimeout(30 * time.Second) {
    log.Println("Warning: some jobs did not complete within 30s")
}
```

## Testing with FakeClock

go-cron includes a built-in `FakeClock` for deterministic testing without real time waits:

```go
fakeClock := cron.NewFakeClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
c := cron.New(cron.WithClock(fakeClock))

executed := make(chan struct{}, 1)
c.AddFunc("0 * * * *", func() {
    executed <- struct{}{}
})
c.Start()
defer c.Stop()

fakeClock.BlockUntil(1)       // Wait for scheduler to register timer
fakeClock.Advance(time.Hour)  // Trigger the job instantly

select {
case <-executed:
    // Job ran successfully
case <-time.After(time.Second):
    t.Fatal("job did not execute")
}
```

No wrapper needed — `cron.NewFakeClock` returns a type that satisfies the `cron.Clock` interface directly.

## Common Options

```go
c := cron.New(
    cron.WithSeconds(),                    // Enable seconds field
    cron.WithLocation(time.UTC),           // Default timezone
    cron.WithContext(parentCtx),            // Parent context
    cron.WithCapacity(100),                // Pre-allocate internals
    cron.WithMaxEntries(1000),             // Limit max entries
    cron.WithRunImmediately(),             // Run @every jobs on Start
    cron.WithLogger(cron.NewSlogLogger(slog.Default())),
    cron.WithChain(cron.Recover(logger)),  // Default wrappers
    cron.WithObservability(hooks),         // Metrics hooks
)
```

## Patterns from Production Usage

### Dynamic Job Management (weaviate pattern)

```go
func (m *Manager) RescheduleJob(name, newSpec string, newJob cron.Job) error {
    // Atomic create-or-update — no manual "check then add/update" needed
    _, err := m.cron.UpsertJob(newSpec, newJob, cron.WithName(name))
    return err
}
```

### Graceful Replacement of Long-Running Jobs

```go
func (m *Manager) ReplaceJob(name, spec string, job cron.Job) error {
    // Wait for current execution to finish before replacing
    m.cron.WaitForJobByName(name)
    _, err := m.cron.UpsertJob(spec, job, cron.WithName(name))
    return err
}
```

### Service Integration with Shutdown

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    c := cron.New(cron.WithContext(ctx))

    c.AddFunc("@every 5m", healthCheck, cron.WithName("health-check"))
    c.AddFunc("0 * * * *", syncData, cron.WithName("hourly-sync"))

    c.Start()

    <-ctx.Done()
    if !c.StopWithTimeout(30 * time.Second) {
        log.Println("Warning: jobs did not complete within 30s")
    }
}
```

### Context-Aware Long-Running Job

```go
c.AddJob("@every 1m", cron.FuncJobWithContext(func(ctx context.Context) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Job canceled, cleaning up")
            return
        case <-ticker.C:
            if err := processNextItem(ctx); err != nil {
                log.Printf("Error: %v", err)
            }
        }
    }
}), cron.WithName("item-processor"))
```

## Migration from robfig/cron

Drop-in replacement — just change the import:

```go
// Before
import "github.com/robfig/cron/v3"

// After
import cron "github.com/netresearch/go-cron"
```

Key behavior differences:
- **DOM/DOW matching**: Uses AND logic (both must match) instead of OR
- **DST spring-forward**: Jobs in skipped hour run immediately instead of being silently skipped
- **Chain execution**: `Entry.Run()` properly invokes chain wrappers

See the [migration guide](https://github.com/netresearch/go-cron/blob/main/docs/MIGRATION.md) for full details.

# Master Refactoring Prompt — URL Shortener Go Backend

You are refactoring a Go URL-shortener backend (Go 1.23, Supabase, Redis). The full source tree is provided. Apply every fix below in a single pass. Remove **all** code comments (inline `//` and block `/* */`) from every file you touch. Do not add any new comments. Preserve existing behaviour unless a fix explicitly changes it. Output every modified file in full.

---

## 1. Critical Bug — JWT Secret Misconfigured

In `cmd/server/main.go`, `jwtSecret` is set to `os.Getenv("DB_API_URL")`. This passes a database URL as the JWT verification endpoint.

- Change it to `os.Getenv("SUPABASE_URL")` (or `JWT_SECRET` — whichever env var holds the Supabase project URL used for JWKS fetching).
- Add a fatal check if the value is empty, consistent with the existing pattern.

---

## 2. Circular Dependency — Service Imports Handler Layer

`internal/service/analytics_service.go` imports `internal/handler/dto` and `internal/handler/mapper`. Services must never reference the handler/presentation layer.

- Move `NormalizeSummary` logic into the service without referencing dto or mapper.
- Cache the `model.UserAnalyticsSummary` directly — marshal/unmarshal the model struct, not the DTO.
- Remove **all** `handler/dto` and `handler/mapper` imports from every file under `internal/service/`.
- The handler layer is the only place that calls mapper functions.

---

## 3. Duplicate `SaveAnalytics` Implementation

`SaveAnalytics` exists in both `url_repository.go` and `analytics_repository.go` with identical logic.

- Remove `SaveAnalytics` from `url_repository.go` and from the `URLRepository` interface.
- All call sites must use `AnalyticsRepository.SaveAnalytics`.

---

## 4. Dead `URLServiceImpl.SaveAnalytics`

In `internal/service/service.go`, `SaveAnalytics` returns `fmt.Errorf("")` — a permanent empty error.

- Remove the method from the `URLService` interface and its implementation entirely.
- Analytics recording already lives in `AnalyticsService.RecordAnalytics`; no replacement is needed.

---

## 5. Duplicate `RespondJSON`

`handler.RespondJSON` (in `analytics_handler.go`, exported) duplicates `utils.RespondJSON`.

- Delete the one in `analytics_handler.go`.
- Update `router.go`'s health-check and any other call site to use `utils.RespondJSON`.

---

## 6. Duplicate Error Response Types

`dto.ErrorResponse` and `handler.ErrorResponse2` overlap.

- Consolidate into a single type in `dto/` with fields: `Error`, `Code` (omitempty), `Field` (omitempty), `RequestID` (omitempty), `Timestamp` (omitempty).
- Delete `ErrorResponse2` from `analytics_handler.go`.
- Update `AnalyticsHandler.respondError` and all other usage to use the unified type.

---

## 7. Goroutine Context Leak in `HandleRedirect`

```go
go func() {
    if err := h.svc.IncrementClickCount(ctx, shortcode); err != nil {
```

`ctx` is derived from the HTTP request and will be cancelled after the response is sent. The goroutine inherits a dead context.

- Create a detached context: `bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)` with a deferred cancel inside the goroutine.

---

## 8. Hardcoded Security Header Dev Mode

In `middleware/security_headers.go`, `IsDevelopment: true` is hardcoded.

- Change `SecurityHeaders()` to accept a `isDev bool` parameter.
- Wire it from the environment config at startup so production deployments get real HSTS/CSP enforcement.

---

## 9. Hardcoded CORS Origin

In `router.go`, `Access-Control-Allow-Origin` is hardcoded to a single Vercel URL.

- Read allowed origins from an env var (`ALLOWED_ORIGINS`), split by comma.
- Validate the request `Origin` header against the list. Reflect the matched origin back; omit the header entirely if no match.
- Keep the existing Vercel URL as the default value if the env var is unset.

---

## 10. Duplicate Rate Limiting in `HandleShorten`

`HandleShorten` in `url_handler.go` has hand-rolled per-IP/user rate limiting with Redis, while a full `RateLimiter` middleware is already applied globally.

- Remove the hand-rolled rate-limit logic from `HandleShorten` (the `Incr`/`Expire`/count-check block).
- If the shorten endpoint needs a stricter limit, register it with `RateLimiter.CustomMiddleware` on that route in the router.

---

## 11. `max()` Shadows Builtin

`middleware/rate_limiter.go` defines a local `max(a, b int) int`. Go 1.21+ has a builtin `max`.

- Delete the custom function; use the builtin everywhere.

---

## 12. Sliding Window Stub

`checkSlidingWindow` just delegates to `checkFixedWindow`.

- Remove `SlidingWindow` from `RateLimiterConfig`.
- Remove the `checkSlidingWindow` method and the branching code in `checkRateLimit`.
- This avoids misleading consumers into thinking sliding window is supported.

---

## 13. Unused `resp` Variable in `GetUserStats`

In `analytics_repository.go → GetUserStats`, the `resp` byte slices from the today-clicks and yesterday-clicks queries are assigned but never read (only the count matters).

- Assign them to `_` to make intent explicit and satisfy linters.

---

## 14. Cache Invalidation Is a No-Op

`AnalyticsServiceImpl.invalidateUserCaches` only logs patterns — it never deletes anything.

- Add a `Delete(ctx context.Context, key string) error` method to the `Cache` interface and implement it in `RedisCache` (via `c.client.Del`).
- Implement actual deletion in `invalidateUserCaches`. For the dashboard key (which is deterministic), delete it directly. For the pattern-based keys, construct the known variants (e.g., specific limit/day values used by the handlers' defaults) and delete those exact keys.

---

## 15. Centralised Config Package

Create `internal/config/config.go` with a struct that reads all env vars once at startup:

```
Port, Environment, AllowedOrigins, RedisURL, SupabaseURL, SupabaseServiceRole,
JWTSecret, Salt, ShortDomain
```

- Parse and validate at startup in `main.go`.
- Pass the config (or relevant fields) down through constructors.
- Remove scattered `os.Getenv` calls from repository, model, cache, and utils packages.

---

## 16. `model.URL.PopulateShortURL` Reads Env at Runtime

`PopulateShortURL` calls `os.Getenv("SHORT_DOMAIN")` on every invocation.

- Change it to accept `baseDomain string` as a parameter.
- The repository that calls it should receive the domain from config at construction time.

---

## 17. `utils.GenerateCode` Reads Env at Runtime

`GenerateCode` reads `SALT` via `os.Getenv` on every call.

- Add `salt string` as a parameter.
- The service that calls it should pass the salt from config.

---

## 18. `cache/keys.go` Reads Env at Runtime

`SecureKey` reads `SALT` from env on every call.

- Accept `salt string` as a parameter.
- Callers pass the salt from config.

---

## 19. Error Handling Consistency

Standardise the error flow: repository → sentinel/wrapped errors, service → wraps with context, handler → maps to HTTP status.

- In `URLRepositoryImpl.GetURLByShortCode`, return `utils.ErrNotFound` when Supabase returns an empty/not-found response instead of a generic error.
- In service methods, propagate the error unchanged (already wrapped).
- In handler methods, use `errors.Is(err, utils.ErrNotFound)` to decide 404 vs 500.

---

## 20. Replace Bubble Sort

`GetUserTopReferrers` and `getUserDeviceBreakdownFromRaw` use manual bubble sort.

- Replace with `slices.SortFunc` (from `slices` package, available in Go 1.21+).

---

## 21. Consistent Nil-Slice JSON Serialization

Several repository methods return `nil` slices that serialize as `null` in JSON.

- Initialise all slice return values to `[]T{}` (empty, not nil) at the start of each function.
- Ensure every JSON response sends `[]` not `null` for list fields.

---

## 22. Structured Logging Scaffolding

Replace raw `log.Printf` calls across the entire codebase with a structured logger.

- Create `internal/logger/logger.go` that wraps `log/slog` (stdlib, Go 1.21+).
- Expose a package-level `Logger` (an `*slog.Logger`) initialised with a JSON handler for production and a text handler for development.
- Provide a constructor: `func NewLogger(env string) *slog.Logger`.
- Initialise in `main.go` from config and pass through (or set a global).
- Replace every `log.Printf` / `log.Println` / `log.Fatal` with the appropriate `slog.Info`, `slog.Warn`, `slog.Error`, or `slog.Debug` call using structured key-value pairs.

---

## 23. Prometheus Metrics Scaffolding

Set up the foundation for Prometheus monitoring without requiring a running Prometheus instance.

### 23a. Metrics Package

Create `internal/metrics/metrics.go`:

- Import `github.com/prometheus/client_golang/prometheus` and `github.com/prometheus/client_golang/promhttp`.
- Define and register the following collectors:

| Metric Name | Type | Labels | Purpose |
|---|---|---|---|
| `http_requests_total` | CounterVec | `method`, `path`, `status` | Total HTTP requests |
| `http_request_duration_seconds` | HistogramVec | `method`, `path` | Request latency distribution |
| `http_requests_in_flight` | Gauge | — | Currently active requests |
| `url_shortens_total` | Counter | — | Total shorten operations |
| `url_redirects_total` | Counter | — | Total redirect operations |
| `cache_hits_total` | CounterVec | `operation` | Cache hits |
| `cache_misses_total` | CounterVec | `operation` | Cache misses |
| `db_query_duration_seconds` | HistogramVec | `operation`, `table` | Database call latency |
| `rate_limit_exceeded_total` | CounterVec | `tier` | Rate limit rejections |
| `analytics_records_total` | Counter | — | Analytics events recorded |

- Expose a `func Handler() http.Handler` that returns `promhttp.Handler()`.
- Expose a `func Init()` that registers all collectors with `prometheus.DefaultRegisterer`. Make it safe to call multiple times (use `sync.Once`).

### 23b. Metrics Middleware

Create `internal/middleware/metrics.go`:

- Implement `MetricsMiddleware(next http.Handler) http.Handler` that:
  - Increments `http_requests_in_flight` on entry, decrements on exit.
  - Records `http_request_duration_seconds` using a deferred timer.
  - Increments `http_requests_total` with method, normalised path (strip dynamic segments like short codes), and status code.
- Use a `statusRecorder` wrapper around `http.ResponseWriter` to capture the status code.

### 23c. Instrumented Cache Wrapper

Create `internal/cache/instrumented.go`:

- Implement a struct `InstrumentedCache` that wraps a `Cache` interface and increments `cache_hits_total` / `cache_misses_total` on every `Get` call.
- Pass-through all other methods, recording `db_query_duration_seconds` where appropriate.
- Provide a constructor: `func NewInstrumentedCache(inner Cache) Cache`.

### 23d. Instrumented Repository Wrapper (Optional Scaffolding)

Create `internal/repository/instrumented.go`:

- Provide a `InstrumentedURLRepository` struct that wraps `URLRepository`.
- On every method call, start a timer, call the inner method, then observe `db_query_duration_seconds` with the operation name and table.
- Provide a constructor: `func NewInstrumentedURLRepository(inner URLRepository) URLRepository`.
- Do the same for `AnalyticsRepository` → `InstrumentedAnalyticsRepository`.

### 23e. Wire It Up

In `cmd/server/main.go`:

- Call `metrics.Init()` at startup.
- Wrap the cache: `rc = cache.NewInstrumentedCache(rc)`.
- Wrap repositories: `urlRepo = repository.NewInstrumentedURLRepository(urlRepo)`, same for analytics.
- Add `metrics.MetricsMiddleware` to the middleware chain in the router.
- Register `GET /metrics` → `metrics.Handler()` in the router (not behind auth middleware).

In `go.mod`, add:

```
github.com/prometheus/client_golang v1.20.5
```

Run `go mod tidy` after.

### 23f. Key Metric Increment Points

Sprinkle counter increments at these specific locations:

- `url_shortens_total` — increment in `URLServiceImpl.CreateShortURL` after successful save.
- `url_redirects_total` — increment in `URLHandler.HandleRedirect` before the `http.Redirect` call.
- `rate_limit_exceeded_total` — increment in `RateLimiter.handleLimitExceeded` with the identifier tier as label.
- `analytics_records_total` — increment in `AnalyticsServiceImpl.RecordAnalytics` after dispatching.

---

## 24. Health Check Enhancement

Expand the `/api/health` endpoint to report dependency readiness:

- Ping Redis (`cache.Ping(ctx)` — add a `Ping(ctx) error` method to `Cache` interface and `RedisCache`).
- Attempt a lightweight Supabase query (e.g., select 1 from urls limit 0) or use the existing client's ping if available.
- Return:

```json
{
  "status": "ok",
  "redis": "connected",
  "database": "connected",
  "version": "<from build-time ldflags or env>"
}
```

- If any dependency is down, return status `"degraded"` with HTTP 200 (or 503 if you prefer — pick one and be consistent).

---

## 25. Graceful Shutdown Enhancements

In `cmd/server/main.go`, the shutdown block should also:

- Close the Redis connection (`Add a Close() error` method to `Cache` interface and `RedisCache`, calling `c.client.Close()`).
- Log the final state of in-flight metrics if Prometheus is initialised.
- Increase the shutdown timeout from 5s to 15s for production (read from config).

---

## 26. Request ID Middleware

`AnalyticsHandler` generates its own request IDs per-handler. This should be a middleware concern.

- Create `internal/middleware/request_id.go` that:
  - Checks for an incoming `X-Request-ID` header; if missing, generates one (using the existing `generateRequestID` logic or `uuid.New()`).
  - Stores it in the request context.
  - Sets it as a response header.
- Provide `func GetRequestID(ctx context.Context) string`.
- Remove per-handler `requestIDGen` usage; handlers read from context via `middleware.GetRequestID(r.Context())`.

---

## 27. OpenTelemetry Trace Scaffolding (Lightweight)

Set up the scaffolding so traces can be enabled later with zero code changes.

- Create `internal/telemetry/telemetry.go`:
  - Define `func InitTracer(serviceName, env, otlpEndpoint string) (shutdown func(context.Context) error, err error)`.
  - If `otlpEndpoint` is empty, return a no-op shutdown and do nothing (traces disabled).
  - If provided, configure an OTLP gRPC exporter → BatchSpanProcessor → TracerProvider and set it as the global provider.
  - Return the shutdown function for graceful cleanup.
- In `main.go`, call `InitTracer` with values from config. Defer the shutdown function.
- Add `go.opentelemetry.io/otel` and `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` to `go.mod`.
- Create `internal/middleware/tracing.go`:
  - Implement a middleware that starts a span per request using `otel.Tracer("http")`.
  - Attach method, path, status code as span attributes.
  - If the global tracer is a no-op (no endpoint configured), this middleware is effectively free.

---

## 28. Dependency Injection Cleanup

In `cmd/server/main.go`, the construction order is messy and mixes concerns.

- Group into clear phases:
  1. **Config** — load and validate all configuration.
  2. **Logger** — initialise structured logger.
  3. **Telemetry** — initialise tracer (no-op if not configured).
  4. **Metrics** — initialise Prometheus collectors.
  5. **Dependencies** — Redis, Supabase client.
  6. **Repositories** — construct, wrap with instrumentation.
  7. **Services** — construct with repos + cache.
  8. **Handlers** — construct with services.
  9. **Router** — construct with handlers + middleware chain.
  10. **Server** — start and await shutdown signal.

---

## Execution Rules

1. Apply all changes across all affected files.
2. After each file, strip every code comment — both `//` and `/* */`.
3. Do **not** add any new comments.
4. Preserve all existing tests. Update imports as needed.
5. For new dependencies (`prometheus/client_golang`, `go.opentelemetry.io/otel`), add them to `go.mod` and run `go mod tidy`.
6. If two fixes conflict, note the conflict and prefer the option that preserves existing behaviour.
7. Output every modified and every new file in full.
8. At the end, output a summary table: file path, change type (modified/new/deleted), and a one-line description of what changed.

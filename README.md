# URL Shortener

A full-stack URL shortener with click analytics, built with Go (backend) and React (frontend).

**Live:** [https://url-shortener-nu-two-32.vercel.app](https://url-shortener-nu-two-32.vercel.app)
**API:** [https://url-shortener-backend-an0e.onrender.com](https://url-shortener-backend-an0e.onrender.com)

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              Client Browser                                  │
└───────────────────────────────────┬──────────────────────────────────────────┘
                                    │ HTTPS
                    ┌───────────────┴───────────────┐
                    │                               │
          ┌─────────▼──────────┐         ┌─────────▼──────────┐
          │   React Frontend   │         │   Go HTTP Server   │
          │   (Vite + TSX)     │         │   (net/http)       │
          │   Vercel           │──REST──▶│   Render           │
          │                    │         │                    │
          │  Pages:            │         │  Middleware Chain: │
          │  • Home            │         │  1. RateLimiter    │
          │  • Dashboard       │         │  2. RequestID      │
          │  • Login/Signup    │         │  3. Metrics        │
          │                    │         │  4. Tracing        │
          │  State:            │         │  5. CORS           │
          │  • Zustand         │         │  6. SecurityHdrs   │
          │  • React Query     │         │  7. Auth (JWT)     │
          └────────────────────┘         └─────────┬──────────┘
                                                   │
                    ┌──────────────────────────────┤
                    │                              │
          ┌─────────▼──────────┐       ┌──────────▼─────────┐
          │      Supabase      │       │    Redis (Upstash) │
          │  (PostgreSQL)      │       │    Cache Layer     │
          │                    │       │                    │
          │  Tables:           │       │  Keys:             │
          │  • urls            │       │  • short_url:{sc}  │
          │  • analytics       │       │  • user_urls:{id}  │
          │  • daily_analytics │       │  • analytics:{...} │
          │                    │       │  • ratelimit:{...} │
          │  Auth:             │       │                    │
          │  • JWT / JWKS      │       │  TTLs: 15m – 24h   │
          └────────────────────┘       └────────────────────┘
```

### Request Lifecycle

```
Browser                  Go Server                       Supabase / Redis
   │                         │                                  │
   │── POST /api/urls ───────▶│                                  │
   │                         │── check Redis cache ────────────▶│
   │                         │◀─ miss ─────────────────────────│
   │                         │── GenerateCode(url, len, salt)   │
   │                         │── SaveURL ──────────────────────▶│
   │                         │◀─ inserted row ─────────────────│
   │                         │── cache.Set(short_url:...) ─────▶│
   │◀─ 201 {short_url} ──────│                                  │
   │                         │                                  │
   │── GET /{shortcode} ─────▶│                                  │
   │                         │── cache.Get(short_url:{sc}) ────▶│
   │                         │◀─ hit ──────────────────────────│
   │                         │── go IncrementClickCount() ─────▶│  (async)
   │◀─ 302 Location: URL ────│                                  │
```

### Internal Package Structure

```
url-shortener-go-backend/
├── cmd/server/
│   └── main.go                   # DI wiring in 10 phases
│
└── internal/
    ├── config/                   # Centralised env var loading
    ├── logger/                   # slog initialisation
    ├── metrics/                  # Prometheus collectors
    ├── telemetry/                # OpenTelemetry tracer scaffolding
    │
    ├── cache/
    │   ├── cache.go              # Cache interface
    │   ├── redis.go              # RedisCache implementation
    │   ├── keys.go               # Secure SHA-256 key helpers
    │   └── instrumented.go       # Metrics-wrapped cache
    │
    ├── repository/
    │   ├── url_interface.go      # URLRepository interface
    │   ├── analytics_interface.go
    │   ├── url_repository.go     # Supabase URL queries
    │   ├── analytics_repository.go
    │   ├── instrumented.go       # DB latency metrics wrappers
    │   └── supabase_repository.go
    │
    ├── service/
    │   ├── service.go            # URLService + URLServiceImpl
    │   └── analytics_service.go  # AnalyticsService + impl
    │
    ├── handler/
    │   ├── url_handler.go        # HTTP handlers for URL ops
    │   ├── analytics_handler.go  # HTTP handlers for analytics
    │   ├── dto/                  # Request/response types
    │   └── mapper/               # model → DTO conversions
    │
    ├── middleware/
    │   ├── auth.go               # JWT validation + JWKS cache
    │   ├── rate_limiter.go       # Tiered per-user/IP limiting
    │   ├── security_headers.go   # HSTS, CSP, XSS headers
    │   ├── request_id.go         # X-Request-ID propagation
    │   ├── metrics.go            # HTTP metrics middleware
    │   └── tracing.go            # OTel span-per-request
    │
    ├── router/
    │   └── router.go             # Route registration + CORS
    │
    ├── model/                    # Domain structs
    └── utils/                    # Shared helpers
```

---

## Features

**URL Management**
- Shorten any URL (anonymous or authenticated)
- Configurable short code length (6–12 chars)
- SHA-256 + random salt generation, collision retry
- Instant redirect via `GET /{shortcode}`

**Analytics**
- Click counting per URL (async, non-blocking)
- Per-user dashboard: total URLs, total clicks, daily trend
- Top URLs by click count
- Referrer breakdown
- Device type breakdown (desktop / mobile / tablet / unknown)
- 7–365 day configurable trend window

**Authentication**
- Supabase Auth (email/password)
- JWT verification against Supabase JWKS endpoint
- In-memory JWKS cache with 1-hour TTL

**Observability**
- Structured JSON logging via `log/slog`
- Prometheus metrics on `GET /metrics`
- OpenTelemetry tracing (no-op until `OTEL_EXPORTER_OTLP_ENDPOINT` is set)
- Request ID (`X-Request-ID`) propagated through all log lines

**Resilience**
- Redis caching with variable TTLs (15 min – 1 hour)
- Instrumented cache with hit/miss counters
- Tiered rate limiting (anonymous: 20 req/min, authenticated: 100, premium: 500)
- Burst handling with 1.5× multiplier
- Graceful shutdown (15s production, 5s development)

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, TypeScript, Vite, Tailwind CSS |
| State | Zustand, TanStack React Query |
| Charts | Recharts |
| Backend | Go 1.23 |
| Auth | Supabase Auth (JWT / JWKS) |
| Database | Supabase (PostgreSQL via postgrest-go) |
| Cache | Redis (Upstash) |
| Metrics | Prometheus (`client_golang` v1.20.5) |
| Tracing | OpenTelemetry v1.33 (OTLP gRPC) |
| Security | `unrolled/secure` (HSTS, CSP, XSS) |
| Frontend Deploy | Vercel |
| Backend Deploy | Render |

---

## Environment Variables

### Backend (`url-shortener-go-backend/.env`)

| Variable | Required | Description |
|----------|----------|-------------|
| `SUPABASE_URL` | ✅ | Supabase project URL — used to fetch JWKS for JWT verification |
| `DB_API_URL` | ✅ | Supabase REST API URL (PostgREST endpoint) |
| `SERVICE_ROLE` | ✅ | Supabase service-role key (bypasses RLS) |
| `SALT` | ✅ | ≥32-char random secret for short code generation and cache key hashing |
| `REDIS_URL` | ✅ | Redis connection URL (`rediss://:<password>@host:port`) |
| `SHORT_DOMAIN` | ✅ | Public base URL for short links (e.g. `https://your.domain`) |
| `ALLOWED_ORIGINS` | ✅ | Comma-separated allowed CORS origins |
| `PORT` | — | HTTP port (default: `8080`) |
| `ENVIRONMENT` | — | `development` or `production` (default: `production`) |
| `APP_VERSION` | — | Build version string shown in `/api/health` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | — | OTLP gRPC endpoint (tracing disabled if unset) |

### Frontend (`url-shortener-frontend/.env`)

| Variable | Description |
|----------|-------------|
| `VITE_SUPABASE_URL` | Supabase project URL |
| `VITE_SUPABASE_ANON_KEY` | Supabase anonymous key |

---

## API Reference

All endpoints return `application/json`. Protected routes require `Authorization: Bearer <jwt>`.

### URLs

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/urls` | optional | Shorten a URL |
| `GET` | `/api/urls` | ✅ | List authenticated user's URLs |
| `GET` | `/api/urls/{shortcode}` | — | Get URL metadata by short code |
| `GET` | `/{shortcode}` | — | Redirect to original URL |

**`POST /api/urls`**
```json
// Request
{ "url": "https://example.com", "is_public": true, "code_length": 7 }

// Response 201
{
  "id": "uuid",
  "short_code": "aBc1234",
  "short_url": "https://your.domain/aBc1234",
  "created_at": "2026-01-01T00:00:00Z",
  "is_public": true,
  "click_count": 0
}
```

### Analytics (all require auth)

| Method | Path | Query Params | Description |
|--------|------|-------------|-------------|
| `GET` | `/api/analytics/dashboard` | — | Aggregated user summary |
| `GET` | `/api/analytics/urls` | `limit` (1–100, default 10) | Top URLs by clicks |
| `GET` | `/api/analytics/referrers` | `limit` (1–50, default 5) | Top referrers |
| `GET` | `/api/analytics/devices` | — | Device type breakdown |
| `GET` | `/api/analytics/trend` | `days` (1–365, default 7) | Daily click trend |
| `POST` | `/api/analytics/record` | — | Record a click event |

**`GET /api/analytics/dashboard`**
```json
{
  "overview": {
    "total_urls": 42,
    "total_clicks": 1337,
    "clicks_today": 23,
    "clicks_yesterday": 18,
    "average_clicks": 31.8,
    "trend_direction": "up"
  },
  "top_urls": [{ "url_id": "...", "short_code": "aBc1234", "original_url": "...", "click_count": 200, "created_at": "..." }],
  "top_referrers": [{ "referrer": "https://twitter.com", "clicks": 80 }],
  "device_breakdown": [{ "device_type": "mobile", "clicks": 900, "percentage": 67.3 }],
  "daily_trend": [{ "date": "2026-01-01", "clicks": 23 }]
}
```

### System

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/health` | — | Dependency health check |
| `GET` | `/metrics` | — | Prometheus metrics scrape endpoint |

**`GET /api/health`**
```json
{ "status": "ok", "redis": "connected", "database": "connected", "version": "1.2.0" }
```

### Error Responses

All errors follow a consistent shape:

```json
{
  "error": "human-readable message",
  "code": "optional machine code",
  "field": "optional field name",
  "request_id": "req_1234567890_42",
  "timestamp": 1735689600
}
```

---

## Rate Limiting

Limits are applied per user (authenticated) or per IP (anonymous).

| Tier | Limit | Window |
|------|-------|--------|
| Anonymous | 20 req | 1 min |
| Authenticated | 100 req | 1 min |
| Premium | 500 req | 1 min |

Burst traffic up to 1.5× the limit is absorbed for up to 30 seconds (production). Rate limit headers are always included:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1735689660
Retry-After: 42          # only on 429
```

---

## Caching Strategy

| Data | Cache Key Pattern | TTL |
|------|-------------------|-----|
| Short URL lookup | `short_url:{shortcode}` | 1 hour |
| User URL list | `user_urls:{userID}` | 1 hour |
| Analytics dashboard | `analytics_{hash}` (HMAC of userID + date range) | 1 hour |
| Top URLs | `user_top_urls:{userID}:{limit}` | 30 min |
| Daily trend | `user_daily_trend:{userID}:{days}` | 15 min |
| Top referrers | `user_top_referrers:{userID}:{limit}` | 45 min |
| Device breakdown | `user_device_breakdown:{userID}` | 1 hour |

Cache keys for user data are hashed with SHA-256 using the server `SALT` to prevent enumeration.

After a `RecordAnalytics` event, the affected user's cache keys are explicitly deleted (the known variants for each default limit).

---

## Prometheus Metrics

Exposed at `GET /metrics` (not behind auth middleware).

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | `method`, `path`, `status` | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | `method`, `path` | Request latency |
| `http_requests_in_flight` | Gauge | — | Active concurrent requests |
| `url_shortens_total` | Counter | — | Total successful shorten operations |
| `url_redirects_total` | Counter | — | Total redirects served |
| `cache_hits_total` | Counter | `operation` | Redis cache hits |
| `cache_misses_total` | Counter | `operation` | Redis cache misses |
| `db_query_duration_seconds` | Histogram | `operation`, `table` | Supabase query latency |
| `rate_limit_exceeded_total` | Counter | `tier` | Rate limit rejections by tier |
| `analytics_records_total` | Counter | — | Analytics events dispatched |

---

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 20+
- A Supabase project with the schema below
- A Redis instance (Upstash recommended for serverless)

### Backend

```bash
cd url-shortener-go-backend
cp .env.example .env        # fill in values
go run ./cmd/server
```

The server starts on `http://localhost:8080`.

### Frontend

```bash
cd url-shortener-frontend
cp .env.example .env.local  # fill in VITE_SUPABASE_URL and VITE_SUPABASE_ANON_KEY
npm install
npm run dev
```

The dev server starts on `http://localhost:5173`.

---

## Database Schema

Run these in the Supabase SQL editor:

```sql
create table urls (
  id          uuid primary key default gen_random_uuid(),
  user_id     uuid references auth.users(id),
  original_url text not null,
  short_code  text not null unique,
  is_public   boolean not null default true,
  click_count bigint not null default 0,
  created_at  timestamptz not null default now()
);

create table analytics (
  id          uuid primary key default gen_random_uuid(),
  url_id      text not null,
  user_id     uuid references auth.users(id),
  referrer    text,
  device_type text,
  clicked_at  timestamptz not null default now()
);

create table daily_analytics (
  id               uuid primary key default gen_random_uuid(),
  url_id           text not null,
  user_id          uuid references auth.users(id),
  date             date not null,
  click_count      bigint default 0,
  unique_referrers bigint default 0,
  desktop_clicks   bigint default 0,
  mobile_clicks    bigint default 0,
  tablet_clicks    bigint default 0,
  unknown_clicks   bigint default 0,
  created_at       timestamptz default now(),
  updated_at       timestamptz default now()
);

-- RPC: increment click count atomically
create or replace function increment_click_count(sc text)
returns void language sql as $$
  update urls set click_count = click_count + 1 where short_code = sc;
$$;

-- RPC: daily click aggregation
create or replace function get_user_daily_clicks(p_user_id uuid, p_days int)
returns table(date text, clicks bigint) language sql as $$
  select
    to_char(date_trunc('day', clicked_at), 'YYYY-MM-DD') as date,
    count(*) as clicks
  from analytics
  where user_id = p_user_id
    and clicked_at >= now() - (p_days || ' days')::interval
  group by date_trunc('day', clicked_at)
  order by 1;
$$;
```

---

## Deployment

### Backend on Render

1. Connect your repository to Render as a **Web Service**
2. Set **Build Command:** `go build -o server ./cmd/server`
3. Set **Start Command:** `./server`
4. Add all required environment variables in the Render dashboard

### Frontend on Vercel

1. Connect the `url-shortener-frontend` directory
2. Set environment variables `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY`
3. Vercel auto-detects Vite and configures the build

---

## Development Notes

- **SALT** must be ≥ 32 characters. Short code generation will hard-fail at startup if it's missing or too short.
- Setting `ENVIRONMENT=development` enables: text-format logging, relaxed rate limits (1000/min), and security headers in dev mode (no HSTS enforcement).
- The `/metrics` endpoint is intentionally unauthenticated — restrict access at the network/proxy layer in production.
- OpenTelemetry tracing is a no-op until `OTEL_EXPORTER_OTLP_ENDPOINT` is set; there is zero overhead when unconfigured.
- Cache invalidation after analytics recording only deletes the known default-limit variants. If you call analytics endpoints with non-default limits, those cache entries will expire naturally per their TTL.
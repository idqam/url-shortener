# URL Shortener

A full-stack URL shortener with click analytics, built with Go (backend) and React (frontend).

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

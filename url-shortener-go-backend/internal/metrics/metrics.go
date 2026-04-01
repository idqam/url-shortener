package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency distribution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Currently active requests",
	})

	URLShortensTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "url_shortens_total",
		Help: "Total shorten operations",
	})

	URLRedirectsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "url_redirects_total",
		Help: "Total redirect operations",
	})

	CacheHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Cache hits",
		},
		[]string{"operation"},
	)

	CacheMissesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Cache misses",
		},
		[]string{"operation"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database call latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	RateLimitExceededTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Rate limit rejections",
		},
		[]string{"tier"},
	)

	AnalyticsRecordsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_records_total",
		Help: "Analytics events recorded",
	})
)

var once sync.Once

func Init() {
	once.Do(func() {
		prometheus.MustRegister(
			HTTPRequestsTotal,
			HTTPRequestDuration,
			HTTPRequestsInFlight,
			URLShortensTotal,
			URLRedirectsTotal,
			CacheHitsTotal,
			CacheMissesTotal,
			DBQueryDuration,
			RateLimitExceededTotal,
			AnalyticsRecordsTotal,
		)
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}

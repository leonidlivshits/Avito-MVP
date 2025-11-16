package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pr_reviewer_http_requests_total",
			Help: "Total number of HTTP requests processed, labeled by method, path and status.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pr_reviewer_http_request_duration_seconds",
			Help:    "HTTP request durations in seconds.",
			Buckets: []float64{0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.3, 0.5, 1.0, 2.5, 5.0},
		},
		[]string{"method", "path"},
	)

	HTTPInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pr_reviewer_http_requests_in_flight",
			Help: "A gauge of requests currently being served.",
		},
	)
)


func RegisterCollectors() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(HTTPInFlight)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func InstrumentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		HTTPInFlight.Inc()
		sr := &statusRecorder{ResponseWriter: w, status: 0}

		next.ServeHTTP(sr, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method
		status := strconv.Itoa(sr.status)

		HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPInFlight.Dec()
	})
}

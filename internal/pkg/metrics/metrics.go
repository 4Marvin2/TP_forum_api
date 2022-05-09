package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	Hits    *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}

func RegisterMetrics(r *mux.Router) *PromMetrics {
	var metrics PromMetrics

	metrics.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "timings",
		},
		[]string{"status", "path", "method"},
	)

	prometheus.MustRegister(metrics.Hits, metrics.Timings)

	r.Handle("/metrics", promhttp.Handler())

	return &metrics
}

func Metrics(metrics *PromMetrics) (mw func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqTime := time.Now()
			next.ServeHTTP(w, r)
			respTime := time.Since(reqTime)
			if r.URL.Path != "/metrics" {
				metrics.Hits.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Inc()
				metrics.Timings.WithLabelValues(
					strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Observe(respTime.Seconds())
			}
		})
	}
}

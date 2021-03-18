package v1

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Stopped = false
)

var (
	goLevDistanceResultsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "go_lev",
			Subsystem: "random",
			Name:      "results",
			Help:      "Go levenshtein random results",
			Buckets:   []float64{1, 5, 10, 20, 30, 50, 100, 250, 500},
		},
		[]string{"count", "max"},
	)
)

func init() {
	InitializeEndpoints()
	InitializeMetrics()
}

func InitializeEndpoints() {
	http.HandleFunc("/v1/random", httpGetRandom)
	http.HandleFunc("/v1/distance", httpGetDistance)
}

func InitializeMetrics() {
	prometheus.MustRegister(goLevDistanceResultsHistogram)
}

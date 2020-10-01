package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/sylr/go-lev/api"
)

var (
	version   = "dev"
	goVersion = runtime.Version()
)

var (
	goLevenshteinBuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "go_lev",
			Subsystem: "",
			Name:      "build_info",
			Help:      "Go levenshtein build info",
		},
		[]string{"version"},
	)
)

func init() {
	// Register build info
	prometheus.MustRegister(goLevenshteinBuildInfo)
	goLevenshteinBuildInfo.WithLabelValues(version).Set(1)
}

func main() {
	http.HandleFunc("/version", httpGetVersion)
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func httpGetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", version)
}

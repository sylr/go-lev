package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	instana "github.com/instana/go-sensor"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	_ "github.com/sylr/go-lev/api"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

var (
	version   = "dev"
	goVersion = runtime.Version()
	sensor    = instana.NewSensor("go-lev")
)

var (
	goLevenshteinBuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "go_lev",
			Subsystem: "",
			Name:      "build_info",
			Help:      "Go levenshtein build info",
		},
		[]string{"version", "goversion"},
	)
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  true,
		DisableSorting: false,
		SortingFunc:    tools.SortLogKeys,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)

	// Register build info
	prometheus.MustRegister(goLevenshteinBuildInfo)
}

func main() {
	http.HandleFunc("/version", httpGetVersion)
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func httpGetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", version)
}

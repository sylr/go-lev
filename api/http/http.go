package http

import (
	"fmt"
	"net/http"
	"sync"

	httptools "github.com/sylr/go-lev/api/http/tools"
	v1 "github.com/sylr/go-lev/api/http/v1"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mutex   = sync.RWMutex{}
	Stopped = false
)

func init() {
	InitializeEndpoints()
}

func InitializeEndpoints() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/start", httptools.OpentracingHTTPWrapper("httpGetStart", httpGetStart))
	http.HandleFunc("/stop", httptools.OpentracingHTTPWrapper("httpGetStop", httpGetStop))
	http.HandleFunc("/ready", httptools.OpentracingHTTPWrapper("httpGetReady", httpGetReady))
	http.HandleFunc("/ping", httptools.OpentracingHTTPWrapper("httpGetPing", httpGetPing))
}

func httpGetStart(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	Stopped = false
	v1.Stopped = false
}

func httpGetStop(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	Stopped = true
	v1.Stopped = true
}

func httpGetReady(w http.ResponseWriter, r *http.Request) {
	if !stopped(w, r) {
		fmt.Fprint(w, "OK")
	}
}

func httpGetPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func stopped(w http.ResponseWriter, r *http.Request) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	if Stopped {
		http.NotFound(w, r)
		return true
	}

	return false
}

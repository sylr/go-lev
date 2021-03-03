package http

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	http.HandleFunc("/start", httpGetStart)
	http.HandleFunc("/stop", httpGetStop)
	http.HandleFunc("/ready", httpGetReady)
	http.HandleFunc("/ping", httpGetPing)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signals:
			mutex.Lock()
			defer mutex.Unlock()
			fmt.Println("SIGTERM received, stop sending 200 OK on /ready.")

			Stopped = true
			signal.Stop(signals)
		}
	}()
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

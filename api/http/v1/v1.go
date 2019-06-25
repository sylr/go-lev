package v1

import (
	"net/http"

	httptools "github.com/sylr/go-lev/api/http/tools"
)

var (
	Stopped = false
)

func init() {
	InitializeEndpoints()
}

// InitializeEndpoints ...
func InitializeEndpoints() {
	http.HandleFunc("/v1/random", httptools.OpentracingHTTPWrapper("httpGetRandom", httpGetRandom))
	http.HandleFunc("/v1/distance", httptools.OpentracingHTTPWrapper("httpGetDistance", httpGetDistance))
}

func stopped(w http.ResponseWriter, r *http.Request) bool {
	if Stopped {
		http.NotFound(w, r)
		return true
	}

	return false
}

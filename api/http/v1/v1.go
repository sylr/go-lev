package v1

import (
	"net/http"
)

var (
	Stopped = false
)

func init() {
	InitializeEndpoints()
}

func InitializeEndpoints() {
	http.HandleFunc("/v1/random", httpGetRandom)
	http.HandleFunc("/v1/distance", httpGetDistance)
}

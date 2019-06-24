package http

import (
	"fmt"
	"net/http"
	"sync"

	v1 "github.com/sylr/go-lev/api/http/v1"

	instana "github.com/instana/go-sensor"
	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mutex   = sync.RWMutex{}
	Stopped = false
)

func init() {
	opentracing.InitGlobalTracer(
		instana.NewTracerWithOptions(
			&instana.Options{
				Service:  "go-lev",
				LogLevel: instana.Info,
			},
		),
	)

	InitializeEndpoints()
}

func InitializeEndpoints() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/start", opentracingWrapper(httpGetStart))
	http.HandleFunc("/stop", opentracingWrapper(httpGetStop))
	http.HandleFunc("/ready", opentracingWrapper(httpGetReady))
	http.HandleFunc("/ping", opentracingWrapper(httpGetPing))
}

func opentracingWrapper(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		parentSpan := opentracing.GlobalTracer().StartSpan("server", ext.RPCServerOption(wireContext))
		parentSpan.SetTag(string(ext.SpanKind), string(ext.SpanKindRPCServerEnum))
		parentSpan.SetTag(string(ext.PeerHostname), r.Host)
		parentSpan.SetTag(string(ext.HTTPUrl), r.URL.Path)
		parentSpan.SetTag(string(ext.HTTPMethod), r.Method)
		parentSpan.SetTag(string(ext.HTTPStatusCode), 200)

		childSpan := opentracing.StartSpan("client", opentracing.ChildOf(parentSpan.Context()))

		f(w, r)

		childSpan.Finish()
		opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(w.Header()))
		parentSpan.Finish()
	}
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

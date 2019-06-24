package v1

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
)

var (
	Stopped = false
)

func init() {
	InitializeEndpoints()
}

func InitializeEndpoints() {
	http.HandleFunc("/v1/random", opentracingWrapper(httpGetRandom))
	http.HandleFunc("/v1/distance", opentracingWrapper(httpGetDistance))
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
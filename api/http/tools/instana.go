package tools

import (
	"net/http"

	instana "github.com/instana/go-sensor"
	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
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
}

// OpentracingHTTPWrapper ...
func OpentracingHTTPWrapper(spanName string, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		parentSpan := opentracing.GlobalTracer().StartSpan("http-server", ext.RPCServerOption(wireContext))
		parentSpan.SetTag(string(ext.SpanKind), string(ext.SpanKindRPCServerEnum))
		parentSpan.SetTag(string(ext.PeerHostname), r.Host)
		parentSpan.SetTag(string(ext.HTTPUrl), r.URL.Path)
		parentSpan.SetTag(string(ext.HTTPMethod), r.Method)
		parentSpan.SetTag(string(ext.HTTPStatusCode), 200)

		childSpan := opentracing.StartSpan(spanName, opentracing.ChildOf(parentSpan.Context()))

		f(w, r)

		childSpan.Finish()
		opentracing.GlobalTracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(w.Header()))
		parentSpan.Finish()
	}
}

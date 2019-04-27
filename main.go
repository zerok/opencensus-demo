package main

import (
	"context"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func main() {
	mux := http.NewServeMux()

	// To export metrics for Prometheus, we have to create a
	// Prometheus exporter, attach it to stats view, and register
	// the exporter with the HTTP request muxer:
	promex, _ := prometheus.NewExporter(prometheus.Options{Namespace: "myapp"})
	view.RegisterExporter(promex)
	mux.Handle("/metrics", promex)

	// For tracing we have to do pretty much the same but using
	// the Jaeger exporter. We also have to tell the exporter
	// where it should send these traces:
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"
	jex, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "myapp",
	})
	if err != nil {
		log.Fatalf("Failed to create Jaeger exporter: %s", err.Error())
	}
	trace.RegisterExporter(jex)

	_, span := trace.StartSpan(context.Background(), "handler-setup")
	// Provide a handler where we can increase the
	// loginFailedTotal metric:
	mux.HandleFunc("/failed-login", loginHandler)
	span.End()

	log.Printf("Starting server on port 8888.")
	if err := http.ListenAndServe("127.0.0.1:8888", mux); err != nil {
		log.Fatalf("Failed to start HTTP server: %s", err.Error())
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Let's put a span around the whole failed-login
	// handler so that we can find out, how long it took:
	ctx, span := trace.StartSpan(ctx, "failed-login")
	span.AddAttributes(trace.StringAttribute("user", "some-user"))
	defer span.End()

	helper(ctx)

	// Increment the stat:
	stats.Record(ctx, loginFailedTotal.M(1))
}

func helper(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "helper")
	span.End()
}

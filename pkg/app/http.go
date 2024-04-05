package app

import (
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

const defaultInternalHTTPPort = 8081

// StartInternalHTTP brings up an HTTP server to listening on the internal
// interface used by kubernetes and service monitors
func StartInternalHTTP() Closeable {
	app := webkit.New()

	// Create our middleware.
	//metricsPlug := middleware.New(middleware.Config{
	//	Recorder: metrics.NewRecorder(metrics.Config{}),
	//})

	//h = middlewarestd.Handler("", metricsPlug, func() {})
	app.Get("/health", func(ctx *webkit.Context) error {
		return ctx.String(200, "Healthy")
	})
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: net.JoinHostPort("", strconv.Itoa(defaultInternalHTTPPort)),
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	return server
}

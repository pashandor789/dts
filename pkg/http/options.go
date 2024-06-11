package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"health"
	"log"
	logchi "log/chi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DefaultPrivateAPIOptions(logger log.Logger) RouterOption {
	return RouterOptions(
		WithRealIP(),
		WithJSONResponse(),
		WithLogging(logchi.New(logger)),
		WithRecover(),
	)
}

func DefaultTechOptions(logger promhttp.Logger, registry prometheus.Gatherer) RouterOption {
	return RouterOptions(
		WithRecover(),
		WithReadinessHandler(),
		WithDebugHandler(),
		WithMetricsHandler(logger, registry),
	)
}

func RouterOptions(options ...RouterOption) func(chi.Router) {
	return func(r chi.Router) {
		for _, option := range options {
			option(r)
		}
	}
}

type RouterOption func(chi.Router)

func WithReadinessHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/readiness", health.Routes())
	}
}

func WithDebugHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/debug", middleware.Profiler())
	}
}

func WithMetricsHandler(logger promhttp.Logger, registry prometheus.Gatherer) RouterOption {
	return func(r chi.Router) {
		r.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: logger}))
	}
}

// WithLogging adds requests logging middleware (aka access log).
func WithLogging(logger log.Logger) RouterOption {
	return func(r chi.Router) {
		r.Use(newStructuredLogger(logger))
	}
}

// WithRealIP adds middleware which helps get real requester's IP, not proxy.
func WithRealIP() RouterOption {
	return func(r chi.Router) {
		r.Use(middleware.RealIP)
	}
}

// WithRecover adds recover middleware, which can catch panics from handlers.
func WithRecover() RouterOption {
	return func(r chi.Router) {
		r.Use(middleware.Recoverer)
	}
}

func WithJSONResponse() RouterOption {
	return func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
	}
}

func newStructuredLogger(logger log.Logger) func(next http.Handler) http.Handler {
	l, ok := logger.(middleware.LogFormatter)
	if !ok {
		return middleware.Logger
	}
	return middleware.RequestLogger(l)
}

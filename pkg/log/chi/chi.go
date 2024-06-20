package chi

import (
	"fmt"
	"net/http"
	"time"

	"dts/pkg/log"
	"github.com/go-chi/chi/v5/middleware"
)

var _ middleware.LogEntry = &LogEntry{}

type LogEntry struct {
	Logger log.Logger
}

func (l *LogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   log.DurationToMs(elapsed),
	})
	l.Logger.Infof("Request complete")
}

func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stacktrace": string(stack),
		"panic":      fmt.Sprintf("%+v", v),
	})
	l.Logger.Errorf("Request panic")
}

var _ middleware.LogFormatter = &Logger{}

type Logger struct {
	log.Logger
}

func New(logger log.Logger) *Logger {
	return &Logger{Logger: logger.WithPrefix(log.LibPrefix)}
}

func (z *Logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	uri := fmt.Sprintf("%s://%s%s", scheme, r.Host, log.MaskSessionID(r.RequestURI))
	fields := log.Fields{
		"http_proto":  r.Proto,
		"http_method": r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"http_scheme": scheme,
		"uri":         uri,
	}

	return &LogEntry{Logger: log.Log(r.Context(), z.WithFields(fields))}
}

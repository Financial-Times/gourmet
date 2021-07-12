package log

import (
	"net/http"
)

type HTTPLogger interface {
	Logger
	WithRequest(r *http.Request) *Logger
	WithStatusCode(code int) *Logger
}

type HTTPStructuredLogger struct {
	*StructuredLogger
}

func NewHTTPStructuredLogger(logger *StructuredLogger) *HTTPStructuredLogger {
	return &HTTPStructuredLogger{
		logger,
	}
}

func (h *HTTPStructuredLogger) WithRequest(r *http.Request) *HTTPStructuredLogger {
	route := r.URL.Path
	query := r.URL.RawQuery
	l := h.WithField("proto", r.Proto).
		WithField("method", r.Method).
		WithField("route", route)

	if len(query) > 0 {
		l.WithField("query", query)
	}
	return h
}

func (h *HTTPStructuredLogger) WithStatusCode(code int) *HTTPStructuredLogger {
	h.WithField("status_code", code)
	return h
}

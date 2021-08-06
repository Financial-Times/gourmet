package http

import (
	"net/http"

	"github.com/Financial-Times/gourmet/log"
)

func WithRequest(r *http.Request) log.Field {
	return func(e *log.Event) {
		e.Put("proto", r.Proto)
		e.Put("method", r.Method)
		e.Put("route", r.URL.Path)

		if len(r.URL.RawQuery) > 0 {
			e.Put("query", r.URL.RawQuery)
		}
	}
}

func WithStatusCode(code int) log.Field {
	return log.WithField("status_code", code)
}

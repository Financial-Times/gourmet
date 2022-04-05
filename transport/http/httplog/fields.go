package httplog

import (
	"net/http"

	"github.com/gotha/gourmet/log"
)

func WithRequestFields(r *http.Request) []log.Field {
	return []log.Field{
		WithMethod(r.Method),
		WithPath(r.URL.Path),
		WithUserAgent(r.UserAgent()),
		WithProtocol(r.Proto),
	}
}

func WithMethod(method string) log.Field {
	return log.WithField("method", method)
}

func WithPath(path string) log.Field {
	return log.WithField("path", path)
}

func WithUserAgent(userAgent string) log.Field {
	return log.WithField("userAgent", userAgent)
}

func WithProtocol(protocol string) log.Field {
	return log.WithField("protocol", protocol)
}

func WithStatusCode(code int) log.Field {
	return log.WithField("statusCode", code)
}

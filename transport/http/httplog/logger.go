package httplog

import (
	"context"
	"net/http"

	"github.com/Financial-Times/gourmet/log"
)

type HTTPRequestLogger struct {
	logger log.Logger
}

func NewHTTPRequestLogger(l log.Logger) *HTTPRequestLogger {
	return &HTTPRequestLogger{
		logger: l,
	}
}

func (l *HTTPRequestLogger) Log(ctx context.Context, code int, req *http.Request) {
	fields := WithRequestFields(req)
	fields = append(fields, WithStatusCode(code))
	l.logger.Info("", fields...)
}

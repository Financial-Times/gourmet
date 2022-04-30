package log

import (
	"context"
	"net/http"

	"github.com/gotha/gourmet/log"
)

func HTTPRequestLogger(logger log.Logger) func(ctx context.Context, code int, req *http.Request) {
	return func(ctx context.Context, code int, req *http.Request) {
		fields := WithRequestFields(req)
		fields = append(fields, WithStatusCode(code))
		logger.Info("", fields...)
	}
}

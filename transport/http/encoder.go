package http

import (
	"context"
	"encoding/json"
	"net/http"
)

// JSONResponseEncoder is the common method to encode all response types to the
// client. Since we're using JSON, there's no reason to provide anything more specific.
// There is also the option to specialize on a per-response (per-method) basis.
func JSONResponseEncoder(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func VoidResponseEncoder(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return nil
}

type StatusCodeResolver func(err interface{}) int

type ErrorJSONEncoder struct {
	statusCodeResolver StatusCodeResolver
}

type ErrorJSONEncoderOption func(enc *ErrorJSONEncoder)

func WithStatusCodeResolver(fn StatusCodeResolver) ErrorJSONEncoderOption {
	return func(enc *ErrorJSONEncoder) {
		enc.statusCodeResolver = fn
	}
}

func NewErrorJSONEncoder(options ...ErrorJSONEncoderOption) *ErrorJSONEncoder {
	enc := &ErrorJSONEncoder{}

	for _, opt := range options {
		opt(enc)
	}

	if enc.statusCodeResolver == nil {
		enc.statusCodeResolver = func(_ interface{}) int {
			return http.StatusInternalServerError
		}
	}

	return enc
}

func (e *ErrorJSONEncoder) Encode(ctx context.Context, err error, w http.ResponseWriter) {
	resp := map[string]string{"error": err.Error()}
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	statusCode := e.statusCodeResolver(err)
	w.WriteHeader(statusCode)
	w.Write(jsonResp)
}

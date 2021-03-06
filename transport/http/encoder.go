package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

// EncodeJSONResponse is the common method to encode all response types to the
// client. Since we're using JSON, there's no reason to provide anything more specific.
// There is also the option to specialize on a per-response (per-method) basis.
func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(endpoint.Failer)
	if ok {
		resErr := e.Failed()
		if resErr != nil {
			// Not a Go kit transport error, but a business-logic error.
			// Provide those as HTTP errors.
			EncodeError(ctx, resErr, w)
			return nil
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// todo - we can try to infer the status code from application specific err
	w.WriteHeader(500)

	if e := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	}); e != nil {
		panic(err)
	}
}

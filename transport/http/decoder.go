package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gotha/gourmet/structloader"
)

type queryParamsDataProvider struct {
	req *http.Request
}

func (p *queryParamsDataProvider) Get(key string) (string, error) {
	if p.req.URL == nil {
		return "", structloader.ErrValNotFound
	}
	val := p.req.URL.Query().Get(key)
	if val == "" {
		return "", structloader.ErrValNotFound
	}
	return val, nil
}

func DecodeRequestFromQueryParameters(req interface{}) httptransport.DecodeRequestFunc {
	return func(_ context.Context, httpReq *http.Request) (interface{}, error) {
		loader := structloader.New(
			&queryParamsDataProvider{httpReq},
			structloader.WithLoaderTagName("query"),
		)

		// clone the empty struct without its values
		newReq := reflect.New(reflect.TypeOf(req).Elem()).Interface()
		err := loader.Load(newReq)
		if err != nil {
			return nil, fmt.Errorf("error decoding request: %w", err)
		}

		return newReq, nil
	}
}

func DecodeRequestFromJSONBody(req interface{}) httptransport.DecodeRequestFunc {
	return func(_ context.Context, httpReq *http.Request) (interface{}, error) {
		newReq := reflect.New(reflect.TypeOf(req).Elem()).Interface()
		body, err := io.ReadAll(httpReq.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		err = json.Unmarshal(body, newReq)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling json body: %w", err)
		}
		return newReq, nil
	}
}

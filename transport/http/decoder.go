package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Financial-Times/gourmet/structloader"
	httptransport "github.com/go-kit/kit/transport/http"
)

type queryParamsDataProvider struct {
	req *http.Request
}

func (p *queryParamsDataProvider) Get(key string) (string, error) {
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

		err := loader.Load(req)
		if err != nil {
			return nil, fmt.Errorf("error decoding request: %w", err)
		}

		return req, nil
	}
}

func DecodeRequestFromJSONBody(req interface{}) httptransport.DecodeRequestFunc {
	return func(_ context.Context, httpReq *http.Request) (interface{}, error) {
		body, err := io.ReadAll(httpReq.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		err = json.Unmarshal(body, req)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling json body: %w", err)
		}
		return req, nil
	}
}

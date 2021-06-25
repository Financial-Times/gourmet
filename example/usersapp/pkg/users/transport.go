package users

import (
	"context"
	"encoding/json"
	"net/http"

	gourmethttp "github.com/Financial-Times/gourmet/transport/http"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(r *mux.Router, s Service, logger log.Logger) *mux.Router {
	e := MakeServerEndpoints(s)

	options := gourmethttp.DefaultServerOptions(logger)

	r.Methods("POST").Path("/user").
		Handler(httptransport.NewServer(
			e.RegisterUserEndpoint,
			decodeRegisterUserRequest,
			gourmethttp.EncodeResponse,
			options...,
		))

	r.Methods("DELETE").Path("/user/{id}").
		Handler(httptransport.NewServer(
			e.RegisterUserEndpoint,
			decodeUnregisterUserRequest,
			gourmethttp.EncodeResponse,
			options...,
		))

	r.Methods("GET").Path("/user/{id}").
		Handler(httptransport.NewServer(
			e.GetUserByIDEndpoint,
			decodeGetUserByIDRequest,
			gourmethttp.EncodeResponse,
			options...,
		))

	r.Methods("GET").Path("/users").
		Handler(httptransport.NewServer(
			e.GetAllUsersEndpoint,
			decodeGetAllUsersRequest,
			gourmethttp.EncodeResponse,
			options...,
		))

	return r
}

func decodeRegisterUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req registerUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeUnregisterUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := gourmethttp.ParseIntPathParam(r, "id", "User ID")
	if err != nil {
		return nil, err
	}
	return unregisterUserRequest{UserID: id}, nil
}

func decodeGetUserByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := gourmethttp.ParseIntPathParam(r, "id", "User ID")
	if err != nil {
		return nil, err
	}
	return getUserByIDRequest{UserID: id}, nil
}

func decodeGetAllUsersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getAllUsersRequest{
		Limit:  gourmethttp.ParseUintQueryParam(r, "limit"),
		Offset: gourmethttp.ParseUintQueryParam(r, "offset"),
	}, nil
}

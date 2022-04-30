package main

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type UserDataRequest struct {
	UserID int
}

type UserDataResponse struct {
	Data string
}

func MakeUserDataEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*UserDataRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request")
		}

		return &UserDataResponse{
			Data: fmt.Sprintf("user data for ID: %d", req.UserID),
		}, nil
	}
}

type UserProfileRequest struct {
	Email string
	Full  bool
}

type UserProfileResponse struct {
	Email string
	Full  bool
}

func MakeUserProfileEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*UserProfileRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request")
		}

		return &UserProfileResponse{
			Email: req.Email,
			Full:  req.Full,
		}, nil
	}
}

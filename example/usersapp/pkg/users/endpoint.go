package users

import (
	"context"

	"github.com/Financial-Times/gourmet/example/usersapp/pkg/storage"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RegisterUserEndpoint   endpoint.Endpoint
	UnregisterUserEndpoint endpoint.Endpoint
	GetAllUsersEndpoint    endpoint.Endpoint
	GetUserByIDEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterUserEndpoint:   MakeRegisterUserEndpoint(s),
		UnregisterUserEndpoint: MakeUnregisterUserEndpoint(s),
		GetAllUsersEndpoint:    MakeGetAllUsersEndpoint(s),
		GetUserByIDEndpoint:    MakeGetUserByIDEndpoint(s),
	}
}

type unregisterUserRequest struct {
	UserID int
}

type unregisterUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r unregisterUserResponse) HTTPError() error { return r.Err }

// UnregisterUser godoc
// @Summary Unregister an existing User
// @Description Unregister an existing User
// @Tags User
// @Param id path string true "User ID"
// @Accept  json
// @Produce  json
// @Router /User/{id} [delete]
func MakeUnregisterUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unregisterUserRequest)
		e := s.UnregisterUser(ctx, req.UserID)
		return unregisterUserResponse{
			Err: e,
		}, nil
	}
}

type registerUserRequest struct {
	User *User
}

type registerUserResponse struct {
	User *User `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r registerUserResponse) HTTPError() error { return r.Err }

// RegisterUser godoc
// @Summary Register a new User
// @Description Register a new User
// @Tags User
// @Param User body User.User true "New User"
// @Accept  json
// @Produce  json
// @Success 200 {object} User.User
// @Router /User [post]
func MakeRegisterUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerUserRequest)
		c, e := s.RegisterUser(ctx, req.User)
		return registerUserResponse{
			User: c,
			Err:  e,
		}, nil
	}
}

type getAllUsersRequest struct {
	Limit  uint
	Offset uint
}

type getAllUsersResponse struct {
	Users []User `json:"users,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func (r getAllUsersResponse) HTTPError() error { return r.Err }

// GetAllUsers godoc
// @Summary List existing Users
// @Description List existing Users
// @Tags User
// @Param limit query int false "User count limit" default(100)
// @Param offset query int false "User count offset" default(0)
// @Accept  json
// @Produce  json
// @Success 200 {array} User.User
// @Router /Users [get]
func MakeGetAllUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAllUsersRequest)
		cc, e := s.GetAllUsers(ctx, &storage.QueryOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
		return getAllUsersResponse{
			Users: cc,
			Err:   e,
		}, nil
	}
}

type getUserByIDRequest struct {
	UserID int
}

type getUserByIDResponse struct {
	User User  `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r getUserByIDResponse) HTTPError() error { return r.Err }

// GetUserByID godoc
// @Summary Get an existing User
// @Description Get an existing User
// @Tags User
// @Param id path string true "User ID"
// @Accept  json
// @Produce  json
// @Router /User/{id} [get]
func MakeGetUserByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserByIDRequest)
		c, e := s.GetUserByID(ctx, req.UserID)
		return getUserByIDResponse{
			User: c,
			Err:  e,
		}, nil
	}
}

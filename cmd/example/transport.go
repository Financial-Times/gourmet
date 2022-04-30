package main

import (
	"context"
	"fmt"
	"io"

	"github.com/gotha/gourmet/cmd"
	"github.com/gotha/gourmet/structloader"
)

func decodeUserDataRequest(loader *structloader.Loader) cmd.RequestDecoder {
	return func(ctx context.Context, r io.Reader) (interface{}, error) {
		type request struct {
			UserID int `flag:"user-id" required:"true"`
		}
		req := &request{}
		err := loader.Load(req)
		if err != nil {
			return nil, fmt.Errorf("could not load request data: %w", err)
		}
		return &UserDataRequest{UserID: req.UserID}, nil
	}
}

func encodeUserDataResponse() cmd.ResponseEncoder {
	return func(ctx context.Context, response interface{}, w io.Writer) error {
		resp, ok := response.(*UserDataResponse)
		if !ok {
			return fmt.Errorf("unexpected response type")
		}

		w.Write([]byte(resp.Data))
		return nil
	}
}

func decodeUserProfileRequest(loader *structloader.Loader) cmd.RequestDecoder {
	return func(ctx context.Context, r io.Reader) (interface{}, error) {
		type request struct {
			Email string `flag:"email" required:"true"`
			Full  bool   `flag:"full" default:"false"`
		}
		req := &request{}
		err := loader.Load(req)
		if err != nil {
			return nil, fmt.Errorf("could not load request data: %w", err)
		}
		return &UserProfileRequest{
			Email: req.Email,
			Full:  req.Full,
		}, nil
	}
}

func encodeUserProfileResponse() cmd.ResponseEncoder {
	return func(ctx context.Context, response interface{}, w io.Writer) error {
		resp, ok := response.(*UserProfileResponse)
		if !ok {
			return fmt.Errorf("unexpected response type")
		}

		output := fmt.Sprintf("displaying profile for %s; full: %v", resp.Email, resp.Full)
		w.Write([]byte(output))
		return nil
	}
}

func RegisterTransport(r *cmd.Router) {
	flagsLoader := cmd.NewArgsDataLoader()

	r.Register("get", func(r *cmd.Router) {
		r.Register("user", func(r *cmd.Router) {
			r.Register("data", nil, cmd.WithCommand(
				cmd.NewCommand(
					MakeUserDataEndpoint(),
					decodeUserDataRequest(flagsLoader),
					encodeUserDataResponse(),
				),
			))
			r.Register("profile", nil, cmd.WithCommand(
				cmd.NewCommand(
					MakeUserProfileEndpoint(),
					decodeUserProfileRequest(flagsLoader),
					encodeUserProfileResponse(),
				),
			))
		})
	})
}

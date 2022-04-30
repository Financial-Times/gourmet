package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-kit/kit/endpoint"
)

type RequestDecoder func(context.Context, io.Reader) (interface{}, error)
type ResponseEncoder func(context.Context, interface{}, io.Writer) error

type commandOptions struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func newDefaultCommandOptions() *commandOptions {
	return &commandOptions{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

type CommandOption func(*commandOptions)

func WithStdin(r io.Reader) CommandOption {
	return func(o *commandOptions) {
		o.stdin = r
	}
}

func WithStdout(w io.Writer) CommandOption {
	return func(o *commandOptions) {
		o.stdout = w
	}
}

func WithStderr(w io.Writer) CommandOption {
	return func(o *commandOptions) {
		o.stderr = w
	}
}

type Command struct {
	dec     RequestDecoder
	enc     ResponseEncoder
	e       endpoint.Endpoint
	options *commandOptions
}

func NewCommand(
	e endpoint.Endpoint,
	dec RequestDecoder,
	enc ResponseEncoder,
	opts ...CommandOption,
) *Command {
	options := newDefaultCommandOptions()
	for _, o := range opts {
		o(options)
	}
	return &Command{
		dec:     dec,
		enc:     enc,
		e:       e,
		options: options,
	}
}

func (c *Command) Execute(ctx context.Context) error {
	req, err := c.dec(ctx, c.options.stdin)
	if err != nil {
		return fmt.Errorf("error decoding request: %w", err)
	}
	resp, err := c.e(ctx, req)
	if err != nil {
		return err
	}
	return c.enc(ctx, resp, c.options.stdout)
}

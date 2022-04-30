package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
)

var (
	ErrNotFound = fmt.Errorf("no route matches this command")
)

type routerOptions struct {
	args    []string
	name    string
	command *Command
}

type RouterOption func(*routerOptions)

func WithName(name string) RouterOption {
	return func(o *routerOptions) {
		o.name = name
	}
}

func WithArgs(args []string) RouterOption {
	return func(o *routerOptions) {
		o.args = args
	}
}

func WithCommand(c *Command) RouterOption {
	return func(o *routerOptions) {
		o.command = c
	}
}

type Router struct {
	args     []string
	name     string
	command  *Command
	children []*Router
}

func NewRouter(opts ...RouterOption) *Router {
	options := &routerOptions{}
	for _, o := range opts {
		o(options)
	}
	if options.args == nil {
		args := []string{"main"}
		num_args := len(os.Args)
		for i := 1; i < num_args; i++ {
			if strings.HasPrefix(os.Args[i], "--") {
				break
			}
			args = append(args, os.Args[i])
		}
		options.args = args
	}
	if options.name == "" {
		options.name = "main"
	}

	return &Router{
		args:     options.args,
		name:     options.name,
		command:  options.command,
		children: make([]*Router, 0),
	}
}

func (r *Router) Register(name string, registrar func(*Router), opts ...RouterOption) {
	opts = append(opts, WithName(name), WithArgs(r.args[1:]))
	router := NewRouter(opts...)
	if registrar != nil {
		registrar(router)
	}
	r.children = append(r.children, router)
}

func (r *Router) match(path string) bool {
	return r.name == path
}

func (r *Router) Run(ctx context.Context) error {
	path := r.args[0]
	if !r.match(path) {
		return ErrNotFound
	}

	// if you have children, but there are no more arguments, obviously you cannot access them, so just return error
	remaining := r.args[1:]
	if len(r.children) > 0 && len(remaining) > 0 {
		for _, c := range r.children {
			if c.match(remaining[0]) {
				return c.Run(ctx)
			}
		}
		return ErrNotFound
	}
	if r.command == nil {
		return ErrNotFound
	}

	return r.command.Execute(ctx)
}

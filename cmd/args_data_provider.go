package cmd

import (
	"os"
	"strings"

	"github.com/gotha/gourmet/structloader"
)

type argsDataProviderOptions struct {
	args []string
}

type ArgsDataProviderOption func(*argsDataProviderOptions)

func WithCustomArgs(args []string) ArgsDataProviderOption {
	return func(o *argsDataProviderOptions) {
		o.args = args
	}
}

type ArgsDataProvider struct {
	args map[string]string
}

func (p *ArgsDataProvider) Get(key string) (string, error) {
	val, exists := p.args[key]
	if !exists {
		return "", structloader.ErrValNotFound
	}
	return val, nil
}

func NewArgsDataProvider(opts ...ArgsDataProviderOption) *ArgsDataProvider {
	options := &argsDataProviderOptions{}
	for _, o := range opts {
		o(options)
	}
	if options.args == nil {
		options.args = os.Args
	}

	args := map[string]string{}
	num_args := len(options.args)
	for i := 0; i < num_args; i++ {
		arg := options.args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimLeft(arg, "--")
			if i+1 < num_args {
				next := options.args[i+1]
				if strings.HasPrefix(next, "--") {
					// if next argument is key we assume this is boolean
					args[key] = "true"
				} else {
					args[key] = next
					i++
				}
			} else {
				// if this is the last argument, assume it is bool and if exists then it is true
				args[key] = "true"
			}
		}
	}

	return &ArgsDataProvider{
		args: args,
	}
}

func NewArgsDataLoader(opts ...ArgsDataProviderOption) *structloader.Loader {
	p := NewArgsDataProvider(opts...)
	return structloader.New(
		p,
		structloader.WithLoaderTagName("flag"),
	)
}

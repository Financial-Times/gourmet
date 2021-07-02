package app

import (
	"context"
)

// Lifecycle -
type Lifecycle struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// App -
type App struct {
	Name       string
	lifecycles []Lifecycle
	terminator Terminator
	// @todo - make timeouts configurable via options
	startTimeout int
	stopTimeout  int
	// @todo - add logger
}

// New - creates new Application
func New(name string, lcs []Lifecycle, opts ...Option) *App {
	options := newDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &App{
		Name:       name,
		lifecycles: lcs,
		terminator: options.terminator,
	}
}

// Start -
func (a *App) Start(ctx context.Context) error {
	for _, l := range a.lifecycles {
		if l.OnStart == nil {
			continue
		}
		err := l.OnStart(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop -
func (a *App) Stop(ctx context.Context) error {
	// @todo - stop in reverse order
	for _, l := range a.lifecycles {
		if l.OnStop == nil {
			continue
		}
		err := l.OnStop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run -
func (a *App) Run(ctx context.Context) error {
	err := a.Start(ctx)
	if err != nil {
		return err
	}

	done := make(chan struct{}, 1)
	a.terminator(done)
	<-done

	err = a.Stop(ctx)
	if err != nil {
		return err
	}
	return nil
}

package app

import (
	"context"
)

type App struct {
	lifecycles []Lifecycle
	terminator Terminator
}

// New - creates new Application
func New(lcs []Lifecycle, opts ...Option) *App {
	options := newDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &App{
		lifecycles: lcs,
		terminator: options.terminator,
	}
}

func (a *App) start(ctx context.Context) error {
	for _, l := range a.lifecycles {
		err := l.OnStart(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) stop(ctx context.Context) error {
	for i := len(a.lifecycles) - 1; i >= 0; i-- {
		err := a.lifecycles[i].OnStop(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	err := a.start(ctx)
	if err != nil {
		return err
	}

	done := make(chan struct{}, 1)
	a.terminator(done)
	<-done

	err = a.stop(ctx)
	if err != nil {
		return err
	}
	return nil
}

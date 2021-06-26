package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Lifecycle -
type Lifecycle struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// App -
type App struct {
	Name       string
	Lifecycles []Lifecycle
	// @todo - make timeouts configurable via options
	startTimeout int
	stopTimeout  int
	// @todo - add logger
}

// New - creates new Application
func New(name string, lcs []Lifecycle) *App {
	return &App{
		Name:       name,
		Lifecycles: lcs,
	}
}

// Start -
func (a *App) Start(ctx context.Context) error {
	for _, l := range a.Lifecycles {
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
	for _, l := range a.Lifecycles {
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

	waitForTermination()

	err = a.Stop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func waitForTermination() {
	ch := make(chan os.Signal, 1)
	// @todo - make configureable what is considered termination signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

package app

import "context"

type Lifecycle interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}

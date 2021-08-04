package http

import (
	"context"
)

type ServerLifecycle struct {
	srv *Server
}

func NewServerLifecycle(srv *Server) *ServerLifecycle {
	return &ServerLifecycle{srv}
}

func (l *ServerLifecycle) OnStart(_ context.Context) error {
	go l.srv.Start()
	return nil
}

func (l *ServerLifecycle) OnStop(_ context.Context) error {
	return l.srv.Close()
}

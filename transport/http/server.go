package http

import (
	"net/http"

	"fmt"
	"time"

	"github.com/gorilla/mux"
)

type serverConfig struct {
	readTimeout     int
	writeTimeout    int
	idleTimeout     int
	handlersTimeout int
	appPort         int
}

// Option - function that can change web server options
type Option func(*serverConfig)

// WithCustomHTTPServerReadTimeout - set custom ReadTimeout for the http server
func WithCustomHTTPServerReadTimeout(timeout int) Option {
	return func(o *serverConfig) {
		o.readTimeout = timeout
	}
}

// WithCustomHTTPServerWriteTimeout - set custom WriteTimeout for the http server
func WithCustomHTTPServerWriteTimeout(timeout int) Option {
	return func(o *serverConfig) {
		o.writeTimeout = timeout
	}
}

// WithCustomHTTPServerIdleTimeout - set custom IdleTimeout for the http server
func WithCustomHTTPServerIdleTimeout(timeout int) Option {
	return func(o *serverConfig) {
		o.idleTimeout = timeout
	}
}

// WithCustomHandlersTimeout - set custom timeout for http handlers
func WithCustomHandlersTimeout(timeout int) Option {
	return func(o *serverConfig) {
		o.handlersTimeout = timeout
	}
}

// WithCustomAppPort - set custom port for the web server to run
func WithCustomAppPort(port int) Option {
	return func(o *serverConfig) {
		o.appPort = port
	}
}

func (o *serverConfig) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newServerDefaultConfig() *serverConfig {
	return &serverConfig{
		appPort:         8080,
		readTimeout:     10,
		writeTimeout:    15,
		idleTimeout:     20,
		handlersTimeout: 14,
	}
}

// RoutesRegistrant - you can pass this function to register routes in the services router
type RoutesRegistrant func(router *mux.Router)

// Server - Default Server with some defaults
type Server struct {
	srv *http.Server
}

// New - create new Server with specified options
func New(registrant RoutesRegistrant, options ...Option) *Server {
	config := newServerDefaultConfig()
	config.ApplyOptions(options...)

	serveMux := http.NewServeMux()

	router := mux.NewRouter()
	if registrant != nil {
		registrant(router)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.appPort),
		Handler:      serveMux,
		ReadTimeout:  time.Duration(config.readTimeout) * time.Second,
		WriteTimeout: time.Duration(config.writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.idleTimeout) * time.Second,
	}
	return &Server{srv}
}

// Start - start server
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Close - stop server
func (s *Server) Close() error {
	return s.srv.Close()
}

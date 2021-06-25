package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Financial-Times/gourmet/apperror"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

type ServerWrapper struct {
	Name           string
	Port           string
	log            *log.Logger
	serverInstance *http.Server
}

type Handler struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
}

func NewHTTPServer(name string, port string, log *log.Logger, handlers ...Handler) *ServerWrapper {
	serveMux := http.NewServeMux()
	for _, h := range handlers {
		serveMux.HandleFunc(h.Path, h.HandlerFunc)
	}
	return &ServerWrapper{
		Name: name,
		Port: port,
		log: log,
		serverInstance: &http.Server{
			Addr:    ":" + port,
			Handler: serveMux,
		}}
}
func (s *ServerWrapper) Start() {
	go func() {
		if err := s.serverInstance.ListenAndServe(); err != nil {
			// log.Printf("HTTP server closing with message: %v", err)
		}
	}()
	// log.Printf("[Start] %s HTTP server on port %s started\n", s.Name, s.Port)
}

func (s *ServerWrapper) Shutdown() {
	// log.Printf("[Shutdown] %s HTTP server is shutting down\n", s.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.serverInstance.Shutdown(ctx); err != nil {
		// log.Fatalf("Unable to stop HTTP server: %v", err)
	}
}

func RequestFinalizerFunc(logger log.Logger) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		route := r.URL.Path
		query := r.URL.RawQuery

		var keyvals []interface{}
		keyvals = append(keyvals, "proto", r.Proto, "method", r.Method, "route", route, "status_code", code)
		if len(query) > 0 {
			keyvals = append(keyvals, "query", query)
		}

		logger.Log(keyvals...)
	}
}

func DefaultServerOptions(logger log.Logger) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(EncodeError),
		httptransport.ServerFinalizer(RequestFinalizerFunc(logger)),
	}
}

func codeFrom(err error) int {
	customErr := err.(apperror.AppError)

	switch customErr.ErrorType {
	case apperror.NotFound:
		return http.StatusNotFound
	// case ErrAlreadyExists, ErrInconsistentIDs:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

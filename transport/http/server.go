package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Financial-Times/gourmet/apperror"
	"github.com/Financial-Times/gourmet/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

type ServerWrapper struct {
	Name           string
	Port           string
	log            log.Logger
	serverInstance *http.Server
}

type Handler struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
}

func NewHTTPServer(name string, port string, log log.Logger, handlers ...Handler) *ServerWrapper {
	serveMux := http.NewServeMux()
	for _, h := range handlers {
		serveMux.HandleFunc(h.Path, h.HandlerFunc)
	}
	return &ServerWrapper{
		Name: name,
		Port: port,
		log:  log,
		serverInstance: &http.Server{
			Addr:    ":" + port,
			Handler: serveMux,
		}}
}
func (s *ServerWrapper) Start() {
	go func() {
		if err := s.serverInstance.ListenAndServe(); err != nil {
			s.log.Info("HTTP server closing with message: %v", err)
		}
	}()
	s.log.Info("[Start] %s HTTP server on port %s started\n", s.Name, s.Port)
}

func (s *ServerWrapper) Shutdown() {
	s.log.Info("[Shutdown] %s HTTP server is shutting down\n", s.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.serverInstance.Shutdown(ctx); err != nil {
		s.log.Error("Unable to stop HTTP server: %v", err)
	}
}

func RequestFinalizerFunc(logger *log.HTTPLogger) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		logger.WithRequest(r).WithStatusCode(code).Info("Response log")
	}
}

func DefaultServerOptions(logger *log.HTTPLogger) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		// TODO: This has to implemented
		// httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
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

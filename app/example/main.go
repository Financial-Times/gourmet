package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gotha/gourmet/app"
)

func main() {
	srv := &http.Server{Addr: ":9999"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Kiro")
	})

	a := app.New([]app.Lifecycle{
		&HTTPServerLifecycle{srv},
	})

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.Run(startCtx); err != nil {
		log.Fatal(err)
	}

}

type HTTPServerLifecycle struct {
	srv *http.Server
}

func (l *HTTPServerLifecycle) OnStart(ctx context.Context) error {
	go l.srv.ListenAndServe()
	return nil
}

func (l *HTTPServerLifecycle) OnStop(ctx context.Context) error {
	return l.srv.Shutdown(ctx)
}

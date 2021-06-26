package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Financial-Times/gourmet/app"
)

func main() {
	srv := &http.Server{Addr: ":9999"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Kiro")
	})

	a := app.New("myApp", []app.Lifecycle{
		{
			OnStart: func(_ context.Context) error {
				go srv.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		},
	})

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.Run(startCtx); err != nil {
		log.Fatal(err)
	}

}

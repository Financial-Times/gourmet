package main

import (
	"net/http"

	ghttp "github.com/gotha/gourmet/transport/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi there"))
	}))
	srv := ghttp.New(mux,
		ghttp.WithCustomAppPort(3000),
	)

	srv.Start()

}

package main

import (
	"flag"
	"fmt"
	"github.com/Financial-Times/gourmet/example/usersapp/pkg/storage"
	"github.com/Financial-Times/gourmet/example/usersapp/pkg/users"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Users System API
// @version 1.0
// @description Demo service demonstrating Go-Kit.
// @termsOfService http://swagger.io/terms/

// @contact.name Tsvetan Dimitrov
// @contact.email tsvetan.dimitrov23@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {

	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	db, err := storage.NewDB("reservations")
	if err != nil {
		panic(err)
	}

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	r := mux.NewRouter()

	r = initUsersHandler(r, db, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()

	logger.Log("exit", <-errs)
}

func initUsersHandler(router *mux.Router, db *storage.Persistence, logger log.Logger) *mux.Router {
	r := users.NewUserRepository(*db)
	s := users.NewUserService(r)
	return users.MakeHTTPHandler(router, s, logger)
}

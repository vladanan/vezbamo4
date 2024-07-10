package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/vladanan/vezbamo4/src/routes"
)

func main() {

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)

	routes.RouterSite(r)
	routes.RouterTests(r)
	routes.RouterAssignments(r)
	routes.RouterUsers(r)
	routes.RouterAPI(r)
	routes.RouterI18n(r)

	routes.ServeStatic(r, "/static/")

	slog.Info("Main done")

	server := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:10000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// var err = srv.ListenAndServe("0.0.0.0:10000", r)
	var err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("error starting server", "error", err.Error())
		os.Exit(1)
	}

}

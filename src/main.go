package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	elr "github.com/vladanan/vezbamo4/src/errorlogres"
	"github.com/vladanan/vezbamo4/src/routes"
)

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)

	routes.RouterSite(r)
	routes.RouterTests(r)
	routes.RouterAssignments(r)
	routes.RouterUsers(r)
	routes.RouterAPI(r)
	routes.RouterI18n(r)

	routes.ServeStatic(r, "/static/")

	log.Print(elr.Green + "Main done" + elr.Reset)

	// go gamesForLearningChannelsAndLogs()

	server := &http.Server{
		Addr: "0.0.0.0:10000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      r,
	}

	// Run our server in a go routine so it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("server closed")
			} else {
				log.Println("error starting server, error:", err.Error())
				// os.Exit(1)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	// var err = srv.ListenAndServe("0.0.0.0:10000", r)
	// var err = server.ListenAndServe()
	// if errors.Is(err, http.ErrServerClosed) {
	// 	slog.Info("server closed")
	// } else if err != nil {
	// 	slog.Error("error starting server", "error", err.Error())
	// 	os.Exit(1)
	// }

}

// isključio u staticcheck "-U1000" da ne javlja za postojeće a neiskorišćene funkcije
func gamesForLearningChannelsAndLogs() {

	s := "jedan"
	cc := make(chan string)
	go func() {
		log.Println(s)
		dva := " dva"
		time.Sleep(time.Second * 0)
		// slog.Info(dva)
		cc <- s + dva
	}()
	slog.Info(<-cc)

}

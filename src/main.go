package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/vladanan/vezbamo4/src/routes"
)

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(routes.GoTo404)

	routes.RouterSite(r)
	routes.RouterTests(r)
	routes.RouterAssignments(r)
	routes.RouterUsers(r)
	routes.RouterAPI(r)
	routes.RouterI18n(r)

	routes.ServeStatic(r, "/static/")

	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	var Blue = "\033[34m"
	var Magenta = "\033[35m"
	var Cyan = "\033[36m"
	var Gray = "\033[37m"
	var White = "\033[97m"
	// log.SetPrefix("\n")
	fmt.Println(Red + Green + Yellow + Blue + Magenta + Cyan + Gray + White + Reset)
	slog.Info(Green + "Main done" + Reset)

	go gamesForLearningChannelsAndLogs()

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
				slog.Info("server closed")
			} else {
				slog.Error("error starting server", "error", err.Error())
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
	slog.Info("shutting down")
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

func gamesForLearningChannelsAndLogs() {

	s := "nešto"
	cc := make(chan string)
	go func() {
		dva := " i nešto drugo test 3"
		log.Println(s)
		time.Sleep(time.Second * 5)
		slog.Info(dva)
		cc <- s + dva
	}()
	log.Println(<-cc)

	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	var Blue = "\033[34m"
	var Magenta = "\033[35m"
	var Cyan = "\033[36m"
	var Gray = "\033[37m"
	var White = "\033[97m"
	// log.SetPrefix("\n")
	fmt.Println(Red + Green + Yellow + Blue + Magenta + Cyan + Gray + White + Reset)

	// log.SetPrefix("")
	// e := fmt.Errorf("neka greška")
	// lE(e)
	// slog.Info("zadnji")
}

func lE(e error) {
	var Reset = "\033[0m"
	// var Red = "\033[31m"
	var Yellow = "\033[33m"
	// log.SetFlags(log.Ltime | log.Lshortfile)
	// log.SetPrefix(Red)
	fmt.Println(Yellow + e.Error() + Reset)
	// defer func() { log.SetFlags(log.LstdFlags); log.SetPrefix(Reset) }()
}

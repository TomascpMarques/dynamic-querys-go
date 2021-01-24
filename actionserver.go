package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"

	"github.com/gorilla/mux"
)

// DQGPORT - the port for where DynamicQuerysGo is located
var DQGPORT = os.Getenv("ENV_GOACTIONS_PORT")

// DEFAULTDQGPORT - default port for DynamicQuerysGo
const DEFAULTDQGPORT = "8000"

func main() {
	// Checks for port configuration for the service
	if DQGPORT == "" {
		DQGPORT = DEFAULTDQGPORT
	}

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	routerMux := mux.NewRouter()
	routerMux.HandleFunc("/actions", actions.Handler)

	// server setup
	srv := &http.Server{
		Handler:      routerMux,
		Addr:         "localhost:" + DQGPORT,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// setup for checking if server is online
	online := make(chan string, 1)

	// server start message
	actions.DQGLogger.Printf("Starting Server on addres: http://%s:%s/%s \n", "localhost", DQGPORT, "actions")
	// prevent server blocking
	go func() {
		online <- "Online"
		if err := srv.ListenAndServe(); err != nil {
			actions.DQGLogger.Println(err)
			actions.DQGLogger.Println("Exiting...")
		}
	}()

	// check to see if function
	err := FunctionHooksHealthCheck()
	if err != nil {
		actions.DQGLogger.Println(err)
		os.Exit(1)
	}

	<-online
	actions.DQGLogger.Println("DynamicActions Server is now Online")

	// Graceful-Shutdown implementation
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL,
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	actions.DQGLogger.Println("Server shutting down")
	os.Exit(0)
}

// FunctionHooksHealthCheck -
func FunctionHooksHealthCheck() error {
	if len(actions.FuncsStorage) == 0 {
		return errors.New("There are no functions to be used as endpoints")
	}
	return nil
}

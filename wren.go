package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/compiledpanda/wren-server-go/server"
)

func main() {
	// TODO #2 call config.Get() and pass in cfg object to server.Setup()

	// Setup our Server
	srv := server.Setup()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("listening on %s \n", srv.Addr)

	// Listen for SIGINT (Ctrl+C) so we can trigger a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // Block until we receive our signal.

	// Shutdown server with a timeout and exit
	// Note: If other services must be shutdown see gorilla/mux README for an example
	// TODO #2 Allow shutdown timeout to be configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

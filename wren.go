package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/boltdb/bolt"
	"github.com/compiledpanda/wren-server-go/server"
)

func main() {
	// TODO #2 call config.Get() and pass in cfg object to server.Setup()
	cfg := &server.Config{}

	// Open our Database
	db, err := bolt.Open("wren.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cfg.DB = db

	// Setup our Server
	srv := server.Setup(cfg)

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
	_ = srv.Shutdown(ctx) // Ignore error since we are only making a best effort anyway
	log.Println("shutting down")
	os.Exit(0)
}

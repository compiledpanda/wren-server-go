package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/compiledpanda/wren-server-go/route"

	"github.com/gorilla/mux"
)

func main() {
	// Create Router
	r := mux.NewRouter()

	// Add Routes
	r.HandleFunc("/", route.Root).Methods("GET")

	// TODO #3 Add Method Not allowed and Not Found Handlers

	// Create Server
	srv := &http.Server{
		// TODO #2 Allow Addr to be configurable
		Addr: "0.0.0.0:8985",
		// TODO #2 Allow configurable timeouts
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("listening on %s \n", srv.Addr)

	// SIGINT (Ctrl+C) will trigger a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Shutdown with a timeout
	// TODO #2 Allow shutdown timeout to be configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

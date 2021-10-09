package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/compiledpanda/wren-server-go/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// TODO #2 Configure logger between console and json
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// TODO #2 Call config.Get() and pass in cfg object to server.Setup()
	cfg := &server.Config{}

	// Setup our Server
	srv, err := server.Setup(cfg)
	// Make sure we close our database
	defer cfg.DB.Close()
	if err != nil {
		log.Fatal().Stack().Err(err)
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Stack().Err(err).Msg("Error Listening or Serving")
		}
	}()
	log.Info().Msgf("listening on %s", srv.Addr)

	// Listen for SIGINT (Ctrl+C) so we can trigger a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // Block until we receive our signal.

	// Shutdown server with a timeout and exit
	// Note: If other services must be shutdown see gorilla/mux README for an example
	// TODO #2 Allow shutdown timeout to be configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	log.Info().Msg("shutting down")
	_ = srv.Shutdown(ctx) // Ignore error since we are only making a best effort anyway
}

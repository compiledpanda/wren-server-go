package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func routes(cfg *Config) *mux.Router {
	// Create Router
	r := mux.NewRouter()

	// Log each request
	r.Use(loggingMiddleware)

	// Add Routes
	r.HandleFunc("/", getRoot).Methods("GET")
	r.HandleFunc("/v1/", v1GetRoot).Methods("GET")
	r.HandleFunc("/v1/metadata", v1GetMetadata(cfg)).Methods("GET")
	r.HandleFunc("/v1/metadata", v1PutMetadata(cfg)).Methods("PUT")

	// All unmatched routes should result in a 405 Method Not Allowed
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(methodNotAllowed)

	return r
}

func openDB(path string) (db *bolt.DB, err error) {
	db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}

	// Ensure buckets exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(REPOSITORY))
		if err != nil {
			return fmt.Errorf("could not create repository bucket: %v", err)
		}
		return nil
	})
	return
}

func Setup(cfg *Config) (srv *http.Server, err error) {
	// Open & setup our Database
	// TODO #2 Pull bolt db name and options from config
	db, err := openDB("wren.db")
	if err != nil {
		return
	}
	cfg.DB = db

	// Create Server
	srv = &http.Server{
		// TODO #2 Allow Addr to be configurable
		Addr: "0.0.0.0:8985",
		// TODO #2 Allow configurable timeouts
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      routes(cfg),
	}

	return
}

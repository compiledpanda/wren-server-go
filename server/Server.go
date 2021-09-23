package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func Routes(cfg *Config) *mux.Router {
	// Create Router
	r := mux.NewRouter()

	// Add Routes
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/v1/", V1GetRoot).Methods("GET")
	r.HandleFunc("/v1/metadata", V1GetMetadata(cfg)).Methods("GET")

	// All unmatched routes should result in a 405 Method Not Allowed
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(MethodNotAllowed)

	return r
}

func Setup(cfg *Config) (srv *http.Server, err error) {
	// Open & setup our Database
	// TODO #2 Pull bolt db name and options from config
	db, err := bolt.Open("wren.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(REPOSITORY))
		if err != nil {
			return fmt.Errorf("could not create repository bucket: %v", err)
		}
		return nil
	})
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
		Handler:      Routes(cfg),
	}

	return
}

package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	// Create Router
	r := mux.NewRouter()

	// Add Routes
	r.HandleFunc("/", Root).Methods("GET")

	// All unmatched routes should result in a 405 Method Not Allowed
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(MethodNotAllowed)

	return r
}

func Setup() *http.Server {

	// Create Server
	srv := &http.Server{
		// TODO #2 Allow Addr to be configurable
		Addr: "0.0.0.0:8985",
		// TODO #2 Allow configurable timeouts
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      Routes(),
	}

	return srv
}

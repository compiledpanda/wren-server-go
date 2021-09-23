package server

import "net/http"

func V1GetMetadata(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ReturnJSON(w, http.StatusOK, "Testing...")
	}
}

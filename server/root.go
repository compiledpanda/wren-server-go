package server

import "net/http"

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello World\n"))
}

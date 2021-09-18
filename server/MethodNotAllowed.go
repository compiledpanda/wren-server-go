package server

import (
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	JsonEncoder(w).Encode(Error{"METHOD_NOT_ALLOWED", "Method Not Allowed"})
}

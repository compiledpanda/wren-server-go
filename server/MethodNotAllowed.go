package server

import (
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ReturnJSON(w, http.StatusMethodNotAllowed, Error{"METHOD_NOT_ALLOWED", "Method Not Allowed"})
}

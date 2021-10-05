package server

import "net/http"

type serverError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	returnJSON(w, http.StatusMethodNotAllowed, serverError{"METHOD_NOT_ALLOWED", "Method Not Allowed"})
}

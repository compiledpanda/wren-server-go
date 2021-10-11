package server

import "net/http"

func authenticate(r *http.Request, cfg *Config) (userId string, e serverError) {
	auth := r.Header.Get("Authorization")
	// Ensure auth header is present
	if auth == "" {
		return "", serverError{"MISSING_AUTHORIZATION_HEADER", "Authorization header must be sent"}
	}
	// Ensure that we have a single value

	// Ensure that the value is "Bearer <jwt>"

	// Ensure that the jwt is valid

	// Validate the JWT
	return
}

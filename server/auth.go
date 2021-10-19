package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

var validAlgorithms = map[string]bool{
	"RS256": true,
	"ES256": true,
}

func authenticate(r *http.Request, cfg *Config) (userId string, e serverError) {
	auth := r.Header.Get("Authorization")
	// Ensure auth header is present
	if auth == "" {
		return "", serverError{"MISSING_AUTHORIZATION_HEADER", "Authorization header must be sent"}
	}
	// Ensure that we have a single value
	values := strings.Split(auth, ",")
	if len(values) != 1 {
		return "", serverError{"MALFORMED_AUTHORIZATION_HEADER", "Authorization header must contain a single value"}
	}

	// Ensure that the value is "Bearer <jwt>"
	parts := strings.Fields(strings.TrimSpace(values[0]))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", serverError{"MALFORMED_AUTHORIZATION_HEADER", "Authorization header must contain Bearer <jwt>"}
	}

	// Parse and validate the token
	_, err := jwt.ParseWithClaims(parts[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure we are using a valid algorithm
		if !validAlgorithms[token.Method.Alg()] {
			return nil, errors.New("invalid algorithm")
		}
		// Ensure kid is present
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid must be present and a string")
		}
		// Ensure claims parsed correctly
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			return nil, errors.New("invalid claims")
		}

		// Get User ID and retrieve key
		userId = claims.Subject
		var pem []byte
		err := cfg.DB.View(func(tx *bolt.Tx) error {
			// Get key (or nil if it does not exist)
			b := tx.Bucket([]byte(USER_KEY))
			pem = b.Get([]byte(fmt.Sprintf("%s:%s", userId, kid)))
			return nil
		})
		// Handle key retrieval errors
		if err != nil {
			log.Error().Stack().Err(err).Msg("Unable to retrieve key")
			return nil, errors.New("unable to retrieve key")
		}
		if pem == nil {
			return nil, errors.New("kid not found")
		}
		// Return the correct key based on the algorithm
		switch alg := token.Method.Alg(); alg {
		case "RS256":
			return jwt.ParseRSAPublicKeyFromPEM(pem)
		case "ES256":
			return jwt.ParseECPublicKeyFromPEM(pem)
		default:
			return nil, errors.New("unknown alg")
		}
	})

	if err != nil {
		return "", serverError{"UNAUTHORIZED", fmt.Sprintf("JWT Validation Failed: %s", err.Error())}
	}

	return
}

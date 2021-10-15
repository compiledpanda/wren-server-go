package server

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
)

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

	// Ensure that the jwt is valid

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("my_secret_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	// Validate the JWT
	return
}

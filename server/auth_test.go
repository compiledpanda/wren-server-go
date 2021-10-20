package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestAuthenticate(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	_ = pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	tt := []struct {
		header string
		userId string
		e      serverError
	}{
		{"", "", serverError{"MISSING_AUTHORIZATION_HEADER", "Authorization header must be sent"}},
		{"Bearer jwt1, Bearer jwt2", "", serverError{"MALFORMED_AUTHORIZATION_HEADER", "Authorization header must contain a single value"}},
		{"Bearer jwt1 jwt2", "", serverError{"MALFORMED_AUTHORIZATION_HEADER", "Authorization header must contain Bearer <jwt>"}},
		{"token jwt1", "", serverError{"MALFORMED_AUTHORIZATION_HEADER", "Authorization header must contain Bearer <jwt>"}},
		{"<invalid alg>", "", serverError{"UNAUTHORIZED", "JWT Validation Failed: invalid algorithm"}},
		{"<missing kid>", "", serverError{"UNAUTHORIZED", "TODO"}},
		{"<invalid claim>", "", serverError{"UNAUTHORIZED", "TODO"}},
		{"<no public key>", "", serverError{"UNAUTHORIZED", "TODO"}},
		{"<mismatched alg>", "", serverError{"UNAUTHORIZED", "TODO"}},
		{"<valid>", "", serverError{"UNAUTHORIZED", "TODO"}},
	}

	for i, tc := range tt {
		r, err := http.NewRequest("GET", "/", nil)
		// We generate JWTs on the fly to avoid timing issues
		var header string
		switch tc.header {
		case "<invalid alg>":
			token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.RegisteredClaims{
				Subject:   "test",
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			})
			token.Header["kid"] = "test"
			pk, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			if err != nil {
				panic(err)
			}
			tokenString, err := token.SignedString(pk)
			if err != nil {
				t.Fatal(err)
			}
			header = fmt.Sprintf("Bearer %s", tokenString)
		default:
			header = tc.header
		}
		r.Header.Add("Authorization", header)
		if err != nil {
			t.Fatal(err)
		}
		userId, e := authenticate(r, &Config{})
		if userId != tc.userId {
			t.Errorf("userId does not match in test %d: got %v want %v", i, userId, tc.userId)
		}
		if e.Code != tc.e.Code {
			t.Errorf("e.Code does not match in test %d: got %v want %v", i, e.Code, tc.e.Code)
		}
		if e.Description != tc.e.Description {
			t.Errorf("e.Description does not match in test %d: got %v want %v", i, e.Description, tc.e.Description)
		}
	}
}

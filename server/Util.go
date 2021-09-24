package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

func ReturnJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	// Do NOT escape characters (i.e. < and >) as html
	enc.SetEscapeHTML(false)
	// Save the bytes! (Disable indentation)
	enc.SetIndent("", "")
	// Explicitly ignore errors, since they can only be caused by trying to marshal unsupported types and values
	_ = enc.Encode(v)
}

func ReturnBytes(w http.ResponseWriter, statusCode int, b []byte) {
	// Calculate the hash
	hasher := sha256.New()
	_, err := hasher.Write(b)
	if err != nil {
		ReturnJSON(w, http.StatusInternalServerError, Error{"INTERNAL_ERROR", "Unable to calculate digest"})
	}
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Set Headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Digest", "sha-256="+hash)
	w.WriteHeader(statusCode)

	// Write body
	_, err = w.Write(b)
	if err != nil {
		log.Printf("Unable to write response: %v", err)
	}
}

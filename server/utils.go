package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func returnJSON(w http.ResponseWriter, statusCode int, v interface{}) {
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

func returnBytes(w http.ResponseWriter, statusCode int, b []byte) {
	// Calculate the hash
	hash, err := calculateSHA256(b)
	if err != nil {
		// TODO log error
		returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to calculate digest"})
		return
	}

	// Set Headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Digest", "sha-256="+hash)
	w.WriteHeader(statusCode)

	// Write body
	_, err = w.Write(b)
	// If we error there isn't really anything we can do, so just log the error internally
	if err != nil {
		log.Printf("Unable to write response: %v", err)
	}
}

func returnEmpty(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

// Return the Base64 URL encoded (RFC 4648) SHA-256 hash
func calculateSHA256(b []byte) (hash string, err error) {
	hasher := sha256.New()
	_, err = hasher.Write(b)
	if err != nil {
		return
	}
	hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return
}

func parseDigestHeader(digest string) (sha string, err error) {
	// Split into multiple parts
	parts := strings.Split(strings.TrimSpace(digest), ",")

	// Ensure we have 1 digest pair, and that it is sha-256
	if len(parts) != 1 {
		return "", errors.New("too many parts")
	}
	if !strings.HasPrefix(parts[0], "sha-256=") {
		return "", errors.New("not sha-256")
	}

	// Return just the base64 encoded value (strip off psha-256= prefix)
	return parts[0][8:], nil
}

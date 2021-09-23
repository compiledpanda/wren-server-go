package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/octet-stream")
	// Calculate and set Digest
	// TODO handle errors
	hasher := sha256.New()
	hasher.Write(b)
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	w.Header().Set("Digest", "sha-256="+hash)
	w.WriteHeader(statusCode)

	_, _ = w.Write(b) // TODO handle errors
}

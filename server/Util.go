package server

import (
	"encoding/json"
	"net/http"
)

func ReturnJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	enc := json.NewEncoder(w)
	// Do NOT escape characters (i.e. < and >) as html
	enc.SetEscapeHTML(false)
	// Save the bytes! (Disable indentation)
	enc.SetIndent("", "")
	// Explicitly ignore errors, since they can only be caused by trying to marshal unsupported types and values
	_ = enc.Encode(v)
}

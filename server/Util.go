package server

import (
	"encoding/json"
	"io"
)

func JsonEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	// Do NOT escape characters (i.e. < and >) as html
	enc.SetEscapeHTML(false)
	// Save the bytes! (Disable indentation)
	enc.SetIndent("", "")

	return enc
}

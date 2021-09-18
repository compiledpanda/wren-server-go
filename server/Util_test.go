package server

import (
	"bytes"
	"testing"
)

func TestJsonEncoder(t *testing.T) {
	type exampleStruct struct {
		Content   string `json:"c"`
		Multiline bool   `json:"m"`
	}

	buf := new(bytes.Buffer)
	enc := JsonEncoder(buf)

	enc.Encode(exampleStruct{"<>", false})

	expected := `{"c":"<>","m":false}
`
	if buf.String() != expected {
		t.Errorf("Error marshaled incorrectly: got %v want %v",
			buf.String(), expected)
	}
}

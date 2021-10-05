package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnJSON(t *testing.T) {
	type exampleStruct struct {
		Content   string `json:"c"`
		Multiline bool   `json:"m"`
	}
	rr := httptest.NewRecorder()

	returnJSON(rr, http.StatusOK, exampleStruct{"<>", false})

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if contentType := rr.Result().Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v", contentType)
	}

	expected := `{"c":"<>","m":false}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestReturnBytes(t *testing.T) {
	rr := httptest.NewRecorder()

	expected := []byte("Some Bytes!")
	returnBytes(rr, http.StatusOK, expected)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if contentType := rr.Result().Header.Get("Content-Type"); contentType != "application/octet-stream" {
		t.Errorf("handler returned wrong Content-Type: got %v", contentType)
	}

	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expected))
	}
}

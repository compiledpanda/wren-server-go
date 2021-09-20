package server

import (
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

	ReturnJSON(rr, http.StatusOK, exampleStruct{"<>", false})

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"c":"<>","m":false}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

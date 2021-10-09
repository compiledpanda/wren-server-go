package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestError(t *testing.T) {
	str, _ := json.Marshal(serverError{"code", "description"})

	expected := `{"code":"code","description":"description"}`
	if string(str) != expected {
		t.Errorf("Error marshaled incorrectly: got %v want %v",
			string(str), expected)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	rr := test.CallHandler(t, methodNotAllowed, "GET", "/", nil, nil)

	test.VerifyErrorResponse(t, rr, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
}

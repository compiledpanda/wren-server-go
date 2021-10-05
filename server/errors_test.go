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
	rr := test.CallHandler(t, methodNotAllowed, "GET", "/", nil)

	expected := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`
	test.VerifyHeader(t, rr.Result().Header.Get("Content-Type"), "application/json")
	test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusMethodNotAllowed, expected)
}

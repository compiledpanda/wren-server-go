package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestMethodNotAllowed(t *testing.T) {
	rr := test.CallHandler(t, MethodNotAllowed, "GET", "/", nil)

	expected := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`
	test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusMethodNotAllowed, expected)
}

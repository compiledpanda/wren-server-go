package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EMethodNotAllowed(t *testing.T) {
	statusCode, headers, body := test.CallGetEndpoint(t, Routes(&Config{}), "/bogus")

	expected := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`
	test.VerifyHeader(t, headers.Get("Content-Type"), "application/json")
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusMethodNotAllowed, expected)
}

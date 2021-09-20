package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EMethodNotAllowed(t *testing.T) {
	statusCode, body := test.CallGetEndpoint(t, Routes(), "/bogus")

	expected := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusMethodNotAllowed, expected)
}

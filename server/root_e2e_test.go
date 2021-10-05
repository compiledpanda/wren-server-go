package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2ERoot(t *testing.T) {
	statusCode, headers, body := test.CallGetEndpoint(t, routes(&Config{}), "/")

	expected := "Hello World\n"
	test.VerifyHeader(t, headers.Get("Content-Type"), "text/plain; charset=utf-8")
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusOK, expected)
}

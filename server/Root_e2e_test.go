package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2ERoot(t *testing.T) {
	statusCode, body := test.CallGetEndpoint(t, Routes(), "/")

	expected := "Hello World\n"
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusOK, expected)
}

package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2ERoot(t *testing.T) {
	res := test.CallGetEndpoint(t, routes(&Config{}), "/")

	expected := "Hello World\n"
	test.VerifyStringResponse(t, res, http.StatusOK, expected)
}

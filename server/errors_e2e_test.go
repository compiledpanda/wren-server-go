package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EMethodNotAllowed(t *testing.T) {
	res := test.CallGetEndpoint(t, routes(&Config{}), "/bogus")

	expected := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`
	test.VerifyJSONResponse(t, res, http.StatusMethodNotAllowed, expected)
}

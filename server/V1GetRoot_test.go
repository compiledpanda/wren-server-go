package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestV1GetRoot(t *testing.T) {
	rr := test.CallHandler(t, V1GetRoot, "GET", "/v1/", nil)

	expected := `{"status":"ONLINE"}
`
	test.VerifyStringResponse(t, rr, http.StatusOK, expected)
}

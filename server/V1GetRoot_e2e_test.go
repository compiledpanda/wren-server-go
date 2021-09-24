package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EV1GetRoot(t *testing.T) {
	statusCode, headers, body := test.CallGetEndpoint(t, Routes(&Config{}), "/v1/")

	expected := `{"status":"ONLINE"}
`
	test.VerifyHeader(t, headers.Get("Content-Type"), "application/json")
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusOK, expected)
}

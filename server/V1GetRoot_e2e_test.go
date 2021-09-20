package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EV1GetRoot(t *testing.T) {
	statusCode, body := test.CallGetEndpoint(t, Routes(), "/v1/")

	expected := `{"status":"ONLINE"}
`
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusOK, expected)
}

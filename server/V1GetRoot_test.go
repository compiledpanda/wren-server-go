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
	test.VerifyHeader(t, rr.Result().Header.Get("Content-Type"), "application/json")
	test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusOK, expected)
}

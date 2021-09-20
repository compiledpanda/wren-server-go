package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestRoot(t *testing.T) {
	rr := test.CallHandler(t, Root, "GET", "/v1/", nil)

	expected := "Hello World\n"
	test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusOK, expected)
}

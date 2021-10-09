package server

import (
	"net/http"
	"testing"

	"github.com/compiledpanda/wren-server-go/test"
)

func TestRoot(t *testing.T) {
	rr := test.CallHandler(t, getRoot, "GET", "/v1/", nil, nil)

	expected := "Hello World\n"
	test.VerifyRecordedStringResponse(t, rr, http.StatusOK, expected)
}

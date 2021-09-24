package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EV1GetMetadata(t *testing.T) {
	db, err := bolt.Open("../test_data/v1_get_metadata_value.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		t.Errorf("DB Open Error: %v", err)
	}
	statusCode, headers, body := test.CallGetEndpoint(t, Routes(&Config{DB: db}), "/v1/metadata")

	expected := "Some Bytes!"
	test.VerifyStringResponse(t, statusCode, string(body), http.StatusOK, expected)
	test.VerifyHeader(t, headers.Get("Content-Type"), "application/octet-stream")
	test.VerifyDigestHeader(t, headers.Get("Digest"), []byte("Some Bytes!"))
}

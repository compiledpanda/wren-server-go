package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/compiledpanda/wren-server-go/test"
)

func TestV1GetMetadata(t *testing.T) {
	tt := []struct {
		file  string
		value []byte
	}{
		{"../test_data/v1_get_metadata_nil.db", []byte("")},
		{"../test_data/v1_get_metadata_value.db", []byte("Some Bytes!")},
	}

	for _, tc := range tt {
		db, err := bolt.Open(tc.file, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			db.Close()
			t.Errorf("DB Open Error: %v", err)
		}
		rr := test.CallHandler(t, V1GetMetadata(&Config{DB: db}), "GET", "/v1/metadata", nil)
		db.Close()

		test.VerifyHeader(t, rr.Result().Header.Get("Content-Type"), "application/octet-stream")
		test.VerifyDigestHeader(t, rr.Result().Header.Get("Digest"), tc.value)
		test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusOK, string(tc.value))
	}
}

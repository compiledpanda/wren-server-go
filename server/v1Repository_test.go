package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/compiledpanda/wren-server-go/test"
)

func TestV1GetRoot(t *testing.T) {
	rr := test.CallHandler(t, v1GetRoot, "GET", "/v1/", nil, nil)

	expected := `{"status":"ONLINE"}
`
	test.VerifyHeader(t, rr.Result().Header.Get("Content-Type"), "application/json")
	test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusOK, expected)
}

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
		rr := test.CallHandler(t, v1GetMetadata(&Config{DB: db}), "GET", "/v1/metadata", nil, nil)
		db.Close()

		test.VerifyHeader(t, rr.Result().Header.Get("Content-Type"), "application/octet-stream")
		test.VerifyDigestHeader(t, rr.Result().Header.Get("Digest"), tc.value)
		test.VerifyStringResponse(t, rr.Code, rr.Body.String(), http.StatusOK, string(tc.value))
	}
}

func TestV1PutMetadata(t *testing.T) {
	tt := []struct {
		headers map[string]string
		body    []byte
		status  int
		code    string
	}{
		{nil, []byte(""), http.StatusBadRequest, "MISSING_DIGEST_HEADER"},
		// {map[string]string{"Digest": "sha-256=47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="}, []byte(""), http.StatusCreated, ""},
	}

	for _, tc := range tt {
		db, err := bolt.Open("../test_data/v1_put_metadata_nil.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			db.Close()
			t.Errorf("DB Open Error: %v", err)
		}
		rr := test.CallHandler(t, v1PutMetadata(&Config{DB: db}), "PUT", "/v1/metadata", tc.headers, tc.body)
		db.Close()
		t.Logf("%d %s", rr.Code, rr.Body.String())
		// test.VerifyError(t, rr.Code, rr.Body.Bytes(), tc.status, tc.code)
	}
}

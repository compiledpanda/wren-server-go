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
	test.VerifyRecordedJSONResponse(t, rr, http.StatusOK, expected)
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

		test.VerifyRecordedByteResponse(t, rr, http.StatusOK, tc.value)
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
		{map[string]string{"Digest": "sha-256=47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=,sha-256=47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="}, []byte(""), http.StatusBadRequest, "MALFORMED_DIGEST_HEADER"},
		{map[string]string{"Digest": "sha-512=47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="}, []byte(""), http.StatusBadRequest, "MALFORMED_DIGEST_HEADER"},
		{map[string]string{"Digest": "sha-256=57DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="}, []byte(""), http.StatusBadRequest, "DIGEST_BODY_MISMATCH"},
		{map[string]string{"Digest": "sha-256=47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="}, []byte(""), http.StatusCreated, ""},
		{map[string]string{"Digest": "sha-256=XnxAAXcb7O-PWjwDgip9txg-lJbbyxKaiO04DTUC9ko="}, []byte("metadata!"), http.StatusCreated, ""},
	}

	for _, tc := range tt {
		db, err := bolt.Open("../test_data/v1_put_metadata_nil.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			db.Close()
			t.Errorf("DB Open Error: %v", err)
		}
		rr := test.CallHandler(t, v1PutMetadata(&Config{DB: db}), "PUT", "/v1/metadata", tc.headers, tc.body)
		db.Close()

		if rr.Code == http.StatusCreated {
			if rr.Body.Len() != 0 {
				t.Error("Body is not empty")
			}
		} else {
			test.VerifyErrorResponse(t, rr, tc.status, tc.code)
		}
	}
}

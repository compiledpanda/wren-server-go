package server

import (
	"net/http"
	"testing"

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
	db, err := openDB("../test_data/TestV1GetMetadata.db")
	if err != nil {
		db.Close()
		t.Fatalf("DB Open Error: %v", err)
	}
	defer db.Close()

	tt := []struct {
		value []byte
	}{
		{[]byte("")},
		{[]byte("Some Bytes!")},
	}

	for _, tc := range tt {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(REPOSITORY))
			return b.Put([]byte(REPOSITORY_METADATA), tc.value)
		})
		if err != nil {
			t.Fatalf("DB Put Error: %v", err)
		}
		rr := test.CallHandler(t, v1GetMetadata(&Config{DB: db}), "GET", "/v1/metadata", nil, nil)
		test.VerifyRecordedByteResponse(t, rr, http.StatusOK, tc.value)
	}
}

func TestV1PutMetadata(t *testing.T) {
	db, err := openDB("../test_data/TestV1PutMetadata.db")
	if err != nil {
		db.Close()
		t.Fatalf("DB Open Error: %v", err)
	}
	defer db.Close()

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
		rr := test.CallHandler(t, v1PutMetadata(&Config{DB: db}), "PUT", "/v1/metadata", tc.headers, tc.body)

		if rr.Code == http.StatusCreated {
			if rr.Body.Len() != 0 {
				t.Error("Body is not empty")
			}
		} else {
			test.VerifyErrorResponse(t, rr, tc.status, tc.code)
		}
	}
}

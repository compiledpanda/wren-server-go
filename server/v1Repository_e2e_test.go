package server

import (
	"net/http"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/compiledpanda/wren-server-go/test"
)

func TestE2EV1GetRoot(t *testing.T) {
	res := test.CallGetEndpoint(t, routes(&Config{}), "/v1/")

	expected := `{"status":"ONLINE"}
`
	test.VerifyJSONResponse(t, res, http.StatusOK, expected)
}

func TestE2EV1GetMetadata(t *testing.T) {
	db, err := openDB("../test_data/TestE2EV1GetMetadata.db")
	if err != nil {
		t.Fatalf("DB Open Error: %v", err)
		db.Close()
	}
	defer db.Close()

	expected := []byte("Some Bytes!")
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(REPOSITORY))
		return b.Put([]byte(REPOSITORY_METADATA), expected)
	})
	if err != nil {
		t.Fatalf("DB Put Error: %v", err)
	}

	res := test.CallGetEndpoint(t, routes(&Config{DB: db}), "/v1/metadata")

	test.VerifyByteResponse(t, res, http.StatusOK, expected)
}

func TestE2EV1PutMetadata(t *testing.T) {
	db, err := openDB("../test_data/TestE2EV1PutMetadata.db")
	if err != nil {
		t.Fatalf("DB Open Error: %v", err)
		db.Close()
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(REPOSITORY))
		return b.Delete([]byte(REPOSITORY_METADATA))
	})
	if err != nil {
		t.Fatalf("DB Delete Error: %v", err)
	}

	expected := []byte("Some Bytes!")

	res := test.CallGetEndpoint(t, routes(&Config{DB: db}), "/v1/metadata")
	test.VerifyByteResponse(t, res, http.StatusOK, nil)

	res = test.CallPutBytesEndpoint(t, routes(&Config{DB: db}), "/v1/metadata", expected)
	test.VerifyJSONResponse(t, res, http.StatusCreated, "")

	res = test.CallGetEndpoint(t, routes(&Config{DB: db}), "/v1/metadata")
	test.VerifyByteResponse(t, res, http.StatusOK, expected)
}

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
)

const REPOSITORY = "repository"
const REPOSITORY_METADATA = "metadata"

type repositoryStatus struct {
	Status string `json:"status"`
}

// GET /v1/
func v1GetRoot(w http.ResponseWriter, r *http.Request) {
	returnJSON(w, http.StatusOK, repositoryStatus{"ONLINE"})
}

// GET /v1/metadata
func v1GetMetadata(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO #15 Authenticate
		// TODO #16 Authorize

		// Get metadata
		var metadata []byte
		err := cfg.DB.View(func(tx *bolt.Tx) error {
			// Get metadata (or nil if it does not exist)
			b := tx.Bucket([]byte(REPOSITORY))
			metadata = b.Get([]byte(REPOSITORY_METADATA))
			return nil
		})
		if err != nil {
			// TODO #18 log error
			returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to retrieve metadata"})
			return
		}

		// Return 200
		returnBytes(w, http.StatusOK, metadata)
	}
}

// PUT /v1/metadata
func v1PutMetadata(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO #15 Authenticate
		// TODO #16 Authorize

		// Ensure Digest header is set
		digest := r.Header.Get("Digest")
		if digest == "" {
			returnJSON(w, http.StatusBadRequest, serverError{"MISSING_DIGEST_HEADER", "Digest header not found"})
			return
		}

		// Parse and validate Digest header
		sha, err := parseDigestHeader(digest)
		if err != nil {
			returnJSON(w, http.StatusBadRequest, serverError{"MALFORMED_DIGEST_HEADER", err.Error()})
			return
		}

		// Ensure body matches digest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO #18 log error
			returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to read body"})
			return
		}
		hash, err := calculateSHA256(body)
		if err != nil {
			// TODO #18 log error
			returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to calculate digest"})
			return
		}
		if sha != hash {
			returnJSON(w, http.StatusBadRequest, serverError{"DIGEST_BODY_MISMATCH",
				fmt.Sprintf("Body does not match digest sha. Digest header: %s Hashed body: %s", sha, hash)})
			return
		}

		// Save body as metadata
		err = cfg.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(REPOSITORY))
			return b.Put([]byte(REPOSITORY_METADATA), body)
		})
		if err != nil {
			// TODO #18 log error
			returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to save metadata"})
			return
		}

		// Return 201
		returnEmpty(w, http.StatusCreated)
	}
}

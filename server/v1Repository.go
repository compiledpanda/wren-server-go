package server

import (
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

		var metadata []byte
		// Start a read-only transaction
		err := cfg.DB.View(func(tx *bolt.Tx) error {
			// Get metadata (or nil if it does not exist)
			b := tx.Bucket([]byte(REPOSITORY))
			metadata = b.Get([]byte(REPOSITORY_METADATA))
			return nil
		})

		if err != nil {
			returnJSON(w, http.StatusInternalServerError, serverError{"INTERNAL_ERROR", "Unable to retrieve metadata"})
		}

		returnBytes(w, http.StatusOK, metadata)
	}
}

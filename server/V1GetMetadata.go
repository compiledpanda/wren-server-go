package server

import (
	"net/http"

	"github.com/boltdb/bolt"
)

func V1GetMetadata(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
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
			ReturnJSON(w, http.StatusInternalServerError, Error{"INTERNAL_ERROR", "Unable to retrieve metadata"})
		}

		ReturnBytes(w, http.StatusOK, metadata)
	}
}

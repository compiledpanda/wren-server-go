package server

import "github.com/boltdb/bolt"

type Config struct {
	DB                 *bolt.DB
	Bootstrap          bool
	BootstrapUserId    string
	BootstrapKeyId     string
	BootstrapPublicKey string
}

package server

import "github.com/boltdb/bolt"

type Config struct {
	DB *bolt.DB
}

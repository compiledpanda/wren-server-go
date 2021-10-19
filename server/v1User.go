package server

import "fmt"

const USER = "user"
const USER_KEY = "user_key"

func userPublicKeyKey(userId string, kid string) []byte {
	return []byte(fmt.Sprintf("%s:%s", userId, kid))
}

func userPublicKeyMetadata(userId string, kid string) []byte {
	return []byte(fmt.Sprintf("%s:%s:metadata", userId, kid))
}

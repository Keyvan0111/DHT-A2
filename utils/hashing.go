package utils

import (
	"crypto/sha256"
)

func ConsistentHash (key string) []byte {

	newHash :=sha256.New()
	hashBytes := newHash.Sum(nil)
	truncatedHash := hashBytes[:8]

	return truncatedHash
}

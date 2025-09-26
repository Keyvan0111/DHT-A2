package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const HASHLEN = 8

func ConsistentHash(key string) (string, int) {
	sum := sha256.Sum256([]byte(key))

	// full hash as hex string
	hashHex := hex.EncodeToString(sum[:])

	// convert to big.Int
	num := new(big.Int).SetBytes(sum[:])

	// ring size = 2^m
	ringSize := new(big.Int).Lsh(big.NewInt(1), uint(HASHLEN))

	// compute position on ring
	mod := new(big.Int).Mod(num, ringSize)

	return hashHex, int(mod.Int64())
}

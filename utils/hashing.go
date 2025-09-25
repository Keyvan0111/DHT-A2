package utils

import (
	"crypto/sha256"
	"math"
	"math/big"
)

const HASHLEN int = 8

func ConsistentHash(key string) int {
    hash := sha256.Sum256([]byte(key))

    num := new(big.Int).SetBytes(hash[:])

    ringSize := int64(math.Pow(2, float64(HASHLEN)))

    mod := new(big.Int).Mod(num, big.NewInt(ringSize))

    return int(mod.Int64())
}

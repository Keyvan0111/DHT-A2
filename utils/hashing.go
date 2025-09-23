package utils

import (
    "crypto/sha256"
    "math/big"
)

func ConsistentHash(key string, m int) int {
    hash := sha256.Sum256([]byte(key))

    num := new(big.Int).SetBytes(hash[:])

    ringSize := new(big.Int).Lsh(big.NewInt(1), uint(m)) // 1 << m
    mod := new(big.Int).Mod(num, ringSize)

    return int(mod.Int64())
}

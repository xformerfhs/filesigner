package main

import (
	"crypto/subtle"
	"golang.org/x/crypto/sha3"
)

// getKeyHash calculates the Shake-128 hash of some key bytes.
func getKeyHash(key []byte) []byte {
	hasher := sha3.NewShake128()
	_, _ = hasher.Write([]byte("public"))
	_, _ = hasher.Write(key)
	_, _ = hasher.Write([]byte("key"))
	result1 := make([]byte, 16)
	_, _ = hasher.Read(result1)
	result2 := make([]byte, 16)
	_, _ = hasher.Read(result2)
	result := make([]byte, 16)
	subtle.XORBytes(result, result1, result2)
	return result
}

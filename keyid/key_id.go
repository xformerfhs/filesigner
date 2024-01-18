package keyid

import (
	"crypto/subtle"
	"filesigner/base32encoding"
	"golang.org/x/crypto/sha3"
)

// ******** Private constants ********

var beginFence = []byte{'k', 'e', 'y', 0x5a}
var endFence = []byte{0xa5, 'h', 's', 'h'}

// ******** Public functions ********

// KeyHash calculates the Shake-128 hash of some key bytes.
func KeyHash(key []byte) []byte {
	hasher := sha3.NewShake128()

	_, _ = hasher.Write(beginFence)
	_, _ = hasher.Write(key)
	_, _ = hasher.Write(endFence)

	rawResult := make([]byte, 32)
	_, _ = hasher.Read(rawResult)

	result := make([]byte, 16)
	subtle.XORBytes(result, rawResult[:16], rawResult[16:])

	return result
}

// KeyId returns the key id of some key bytes.
func KeyId(key []byte) string {
	return base32encoding.EncodeKey(KeyHash(key))
}

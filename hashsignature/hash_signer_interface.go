package hashsignature

import (
	"errors"
)

// IsDestroyedErr is returned after the private key of the hash-signer has been destroyed.
var IsDestroyedErr = errors.New("Hash-signer has been destroyed")

// HashSigner is the interface that each HashSigner has to use.
type HashSigner interface {
	GetPublicKey() ([]byte, error)

	SignHash(hashValue []byte) ([]byte, error)

	Destroy()
}

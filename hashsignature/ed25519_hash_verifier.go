package hashsignature

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"strings"
)

// ******** Public types ********

// Ed25519HashVerifier contains the objects necessary for file signature verification
type Ed25519HashVerifier struct {
	publicKey []byte
	options   *ed25519.Options
}

// ******** Private constants ********
const ed25519MinKeyLength = 32
const ed25519MaxKeyLength = 34

// ******** Type creation ********

// NewEd25519HashVerifier creates a new Ed25519HashVerifier.
func NewEd25519HashVerifier(publicKey []byte) (HashVerifier, error) {
	lenKey := len(publicKey)
	if lenKey < ed25519MinKeyLength || lenKey > ed25519MaxKeyLength {
		return nil, fmt.Errorf("bad ed25519 public key length: %d", lenKey)
	}

	result := &Ed25519HashVerifier{
		publicKey: publicKey,
		options:   &ed25519.Options{Hash: crypto.SHA512, Context: fileSignerContext},
	}

	return result, nil
}

// ******** Public functions ********

// VerifyHash verifies the supplied hash with the supplied signature.
func (hv *Ed25519HashVerifier) VerifyHash(hashValue []byte, signature []byte) (bool, error) {
	result, err := hv.doVerifyHash(hashValue, signature)
	if err != nil && strings.Contains(err.Error(), "invalid signature") {
		err = nil
	}

	return result, err
}

// ******** Private functions ********

// doVerifyHash verifies a supplied hash value with as supplied signature.
func (hv *Ed25519HashVerifier) doVerifyHash(hashValue []byte, signature []byte) (bool, error) {
	err := ed25519.VerifyWithOptions(hv.publicKey, hashValue, signature, hv.options)

	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

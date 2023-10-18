package hashsigner

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
)

// ******** Public types ********

// Ec521HashVerifier contains the objects necessary for file curve sec521r1 and hash verification.
type Ec521HashVerifier struct {
	publicKey     *ecdsa.PublicKey
	publicKeyHash []byte
}

// ******** Private constants ********
const p521KeyLength = 158

// ******** Type creation ********

// NewEc521HashVerifier creates a new Ec521HashVerifier.
func NewEc521HashVerifier(publicKey []byte) (HashVerifier, error) {
	lenKey := len(publicKey)
	if lenKey != p521KeyLength {
		return nil, fmt.Errorf("bad ec dsa p521 public key length: %d", lenKey)
	}

	var err error
	result := &Ec521HashVerifier{}

	var pk any
	pk, err = x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("Invalid public key: %v", err)
	}

	var ok bool
	result.publicKey, ok = pk.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Public key is not an ECDSA key")
	}

	return result, nil
}

// ******** Public functions ********

// VerifyHash verifies the supplied hash with the supplied signature.
func (hv *Ec521HashVerifier) VerifyHash(hashValue []byte, signature []byte) (bool, error) {
	return hv.doVerifyHash(hashValue, signature)
}

// ******** Private functions ********

// doVerifyHash verifies a supplied hash value with a supplied signature.
func (hv *Ec521HashVerifier) doVerifyHash(hashValue []byte, signature []byte) (bool, error) {
	return ecdsa.VerifyASN1(hv.publicKey, hashValue, signature), nil
}

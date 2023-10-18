package hashsigner

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
)

// ******** Public types ********

// EcDsa521HashSigner contains the objects necessary for file signing with curve secp521r1 and hash.
type EcDsa521HashSigner struct {
	privateKey *ecdsa.PrivateKey
	isValid    bool
}

// ******** Type creation ********

// NewEcDsa521HashSigner creates a new EcDsa521HashSigner.
func NewEcDsa521HashSigner() (HashSigner, error) {
	var err error

	result := &EcDsa521HashSigner{
		isValid: true,
	}

	result.privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Public functions ********

// GetPublicKey returns a copy of the public key
func (hs *EcDsa521HashSigner) GetPublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(hs.privateKey.Public())
}

// -------- Sign functions --------

// SignHash signs the supplied has value.
func (hs *EcDsa521HashSigner) SignHash(hashValue []byte) ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return hs.doSignHash(hashValue)
}

func (hs *EcDsa521HashSigner) Destroy() {
	if hs.isValid {
		hs.privateKey.D.SetInt64(-1)
		hs.isValid = false
	}
}

// ******** Private functions ********

// checkValidity checks if this Ed25519HashSigner is usable.
func (hs *EcDsa521HashSigner) checkValidity() error {
	if hs.isValid {
		return nil
	} else {
		return IsDestroyedErr
	}
}

// doSignHash signs a supplied hash value.
func (hs *EcDsa521HashSigner) doSignHash(hashValue []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, hs.privateKey, hashValue)
}

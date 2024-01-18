package hashsignature

import (
	"crypto"
	"crypto/ed25519"
	"filesigner/slicehelper"
)

// ******** Public types ********

// Ed25519HashSigner contains the objects necessary for file signing
type Ed25519HashSigner struct {
	signer    ed25519.PrivateKey
	publicKey []byte
	options   *ed25519.Options
	isValid   bool
}

// ******** Type creation ********

// NewEd25519HashSigner creates a new Ed25519HashSigner.
func NewEd25519HashSigner() (HashSigner, error) {
	var err error

	result := &Ed25519HashSigner{
		options: &ed25519.Options{Hash: crypto.SHA512, Context: fileSignerContext},
		isValid: true,
	}

	result.publicKey, result.signer, err = ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ******** Public functions ********

// GetPublicKey returns a copy of the public key
func (hs *Ed25519HashSigner) GetPublicKey() ([]byte, error) {
	return slicehelper.MakeCopy(hs.publicKey), nil
}

// -------- Sign functions --------

// SignHash signs the supplied has value.
func (hs *Ed25519HashSigner) SignHash(hashValue []byte) ([]byte, error) {
	err := hs.checkValidity()
	if err != nil {
		return nil, err
	}

	return hs.doSignHash(hashValue)
}

// Destroy removes the private key from this Ed25519HashSigner, so it can no longer be used.
func (hs *Ed25519HashSigner) Destroy() {
	if hs.isValid {
		slicehelper.ClearInteger(hs.signer)
		hs.signer = nil
		hs.isValid = false
	}
}

// ******** Private functions ********

// checkValidity checks if this Ed25519HashSigner is usable.
func (hs *Ed25519HashSigner) checkValidity() error {
	if hs.isValid {
		return nil
	} else {
		return IsDestroyedErr
	}
}

// doSignHash signs a supplied hash value.
func (hs *Ed25519HashSigner) doSignHash(hashValue []byte) ([]byte, error) {
	return hs.signer.Sign(
		nil,
		hashValue,
		hs.options,
	)
}

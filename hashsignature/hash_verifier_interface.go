package hashsignature

// HashVerifier is the interface that each HashVerifier has to use.
type HashVerifier interface {
	VerifyHash(hashValue []byte, signature []byte) (bool, error)
}

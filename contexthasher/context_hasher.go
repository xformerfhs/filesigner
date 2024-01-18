package contexthasher

import (
	"hash"
)

// ******** Public types ********

// ContextHasher is a hash.Hash with a context.
type ContextHasher struct {
	hasher  hash.Hash
	context []byte
}

// ******** Private constants ********

// separator separates the context from the data.
var separator = []byte{0x33, 0x17, 0xd1, 0xdb, 0xc2, 0xf1}

// ******** Type creation functions ********

// NewContextHasher creates a new context hasher.
func NewContextHasher(hashFunc hash.Hash, contextBytes []byte) hash.Hash {
	result := &ContextHasher{hasher: hashFunc}

	// context points to the bytes of the string. It is not a copy.
	result.context = contextBytes
	result.Reset()

	return result
}

// ******** Public functions ********

// Write writes data to the context hasher.
func (ch *ContextHasher) Write(p []byte) (int, error) {
	return ch.hasher.Write(p)
}

// Sum returns the hash either into b or creates a new byte slice.
func (ch *ContextHasher) Sum(b []byte) []byte {
	return ch.hasher.Sum(b)
}

// Reset resets the context hasher.
func (ch *ContextHasher) Reset() {
	hasher := ch.hasher
	hasher.Reset()
	hasher.Write(ch.context)
	hasher.Write([]byte{byte(len(ch.context))})
	hasher.Write(separator)
}

// Size returns the size of the hash value.
func (ch *ContextHasher) Size() int {
	return ch.hasher.Size()
}

// BlockSize returns the block size of the hash algorithm.
func (ch *ContextHasher) BlockSize() int {
	return ch.hasher.BlockSize()
}

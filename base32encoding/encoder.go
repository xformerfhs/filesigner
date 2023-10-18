package base32encoding

import (
	"strings"
)

// ******** Public functions ********

// EncodeToString encodes a byte slice as a string.
func EncodeToString(b []byte) string {
	emx.Lock()
	defer emx.Unlock()

	return enc.EncodeToString(b)
}

// EncodeToBytes encodes a byte slice as a byte slice.
func EncodeToBytes(b []byte) []byte {
	emx.Lock()
	defer emx.Unlock()

	result := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(result, b)

	return result
}

const keyGroupSize = 4

// EncodeKey encodes a key with groups of letters and numbers.
func EncodeKey(k []byte) string {
	encodedKey := encKey.EncodeToString(k)

	var sb strings.Builder

	f := 0
	t := keyGroupSize
	l := len(encodedKey)

	sb.Grow(l + l/keyGroupSize)

	for {
		if f > 0 {
			sb.WriteByte(keySeparator)
		}

		if t < l {
			sb.WriteString(encodedKey[f:t])
		} else {
			sb.WriteString(encodedKey[f:l])
			break
		}

		f = t
		t += keyGroupSize
	}

	return sb.String()
}

package base32encoding

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
	l := len(encodedKey)
	if l == 0 {
		return ""
	}

	di := 0

	result := make([]byte, l+(l-1)/keyGroupSize+1)
	for l > keyGroupSize {
		result[di] = keySeparator
		di++

		copy(result[di:], encodedKey[:keyGroupSize])

		encodedKey = encodedKey[keyGroupSize:]
		di += keyGroupSize
		l -= keyGroupSize
	}

	result[di] = keySeparator
	di++
	copy(result[di:], encodedKey[:])

	return string(result[1:])
}

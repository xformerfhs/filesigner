package base32encoding

import "strings"

// ******** Public functions ********

// DecodeFromString decodes a string into a byte slice.
func DecodeFromString(s string) ([]byte, error) {
	emx.Lock()
	defer emx.Unlock()

	return enc.DecodeString(s)
}

// DecodeKey decodes a key
func DecodeKey(s string) ([]byte, error) {
	sl := len(s)
	result := make([]byte, sl-(sl-1)/5)

	di := 0
	t := 0

	for {
		t = strings.IndexByte(s, keySeparator)
		if t < 0 {
			copy(result[di:], s)
			break
		}

		copy(result[di:], s[:t])
		di += t
		s = s[t+1:]
	}

	return encKey.DecodeString(string(result))
}

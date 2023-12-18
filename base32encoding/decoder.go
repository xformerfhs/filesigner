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
	result := make([]byte, sl)

	di := 0
	t := 0

	for {
		t = strings.IndexByte(s, keySeparator)
		if t < 0 {
			copy(result[di:], s)
			di += sl
			break
		}

		copy(result[di:], s[:t])
		di += t
		t++
		s = s[t:]
		sl -= t
	}

	return encKey.DecodeString(string(result[:di]))
}

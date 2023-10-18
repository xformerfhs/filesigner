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
	var sb strings.Builder
	sl := len(s)
	if sl > sb.Cap() {
		sb.Grow(sl - sb.Cap())
	}

	f := 0
	t := 0

	for {
		t = strings.IndexByte(s, keySeparator)
		if t < 0 {
			sb.WriteString(s)
			break
		}

		sb.WriteString(s[f:t])
		s = s[t+1:]
	}

	return encKey.DecodeString(sb.String())
}

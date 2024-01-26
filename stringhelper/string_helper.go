package stringhelper

import (
	"strings"
	"unsafe"
)

// UnsafeStringBytes returns a byte slice that points to the bytes of the supplied string.
// No bytes are copied. Attention: This is *unsafe*! Do not change those bytes!
func UnsafeStringBytes(s string) []byte {
	// This is a streamlined version of
	// https://josestg.medium.com/140x-faster-string-to-byte-and-byte-to-string-conversions-with-zero-allocation-in-go-200b4d7105fc .
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeStringFrom returns a string that starts at the specified offset in another string.
// This avoids copying a string.
func UnsafeStringFrom(s string, offset int) string {
	return unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(unsafe.StringData(s)), offset)), len(s)-offset)
}

// HasCaseInsensitivePrefix tests whether the string s begins with prefix in a case-insensitive way.
func HasCaseInsensitivePrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

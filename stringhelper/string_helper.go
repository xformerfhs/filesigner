package stringhelper

import "unsafe"

// UnsafeStringBytes returns a byte slice that points to the bytes of the supplied string.
// No bytes are copied. Attention: This is *unsafe*! Do not change those bytes!
func UnsafeStringBytes(s string) []byte {
	// This is a streamlined version of
	// https://josestg.medium.com/140x-faster-string-to-byte-and-byte-to-string-conversions-with-zero-allocation-in-go-200b4d7105fc .
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

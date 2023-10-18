package base32encoding

import (
	"bytes"
	cryptorand "crypto/rand"
	"math/rand"
	"testing"
)

func TestEncodeKey(t *testing.T) {
	for i := 0; i < 100; i++ {
		sl := rand.Intn(30)
		s := make([]byte, sl)
		_, _ = cryptorand.Read(s)
		es := EncodeKey(s)
		ds, err := DecodeKey(es)
		if err != nil {
			t.Fatalf("Error decoding key '%s': %v", es, err)
		}

		if !bytes.Equal(s, ds) {
			t.Fatalf("Decoding '%s' did not result in '%x', but '%x'", es, s, ds)
		}
	}
}

func BenchmarkEncodeKey(b *testing.B) {
	source := make([]byte, 16)
	_, _ = cryptorand.Read(source)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = EncodeKey(source)
	}
}

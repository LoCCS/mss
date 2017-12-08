package winternitz

import (
	"crypto/rand"
	"testing"
)

func BenchmarkHashToBlocks(b *testing.B) {
	hash := make([]byte, SecurityLevel)

	for i := 0; i < b.N; i++ {
		rand.Read(hash)
		hashToBlocks(hash)
	}
}
